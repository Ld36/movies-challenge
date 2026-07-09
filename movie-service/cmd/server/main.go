package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcadapter "github.com/luizdavid/movies-challenge/movie-service/internal/adapters/grpc"
	mongoadapter "github.com/luizdavid/movies-challenge/movie-service/internal/adapters/mongodb"
	"github.com/luizdavid/movies-challenge/movie-service/internal/application"
	"github.com/luizdavid/movies-challenge/movie-service/internal/config"
	"github.com/luizdavid/movies-challenge/movie-service/internal/database"
	"github.com/luizdavid/movies-challenge/movie-service/internal/logger"
	moviepb "github.com/luizdavid/movies-challenge/movie-service/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	appLogger, err := logger.New(cfg.AppEnv)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	ctx := context.Background()

	mongoClient, err := database.NewMongoClient(ctx, cfg.MongoURI)
	if err != nil {
		appLogger.Fatal("failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			appLogger.Error("failed to disconnect MongoDB", zap.Error(err))
		}
	}()

	collection := mongoClient.
		Database(cfg.MongoDatabase).
		Collection(cfg.MongoCollection)

	movieRepository := mongoadapter.NewMovieRepository(collection)
	movieService := application.NewMovieService(movieRepository)
	movieGRPCServer := grpcadapter.NewMovieGRPCServer(movieService)

	grpcServer := grpc.NewServer()

	moviepb.RegisterMovieServiceServer(grpcServer, movieGRPCServer)

	address := fmt.Sprintf(":%s", cfg.GRPCPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		appLogger.Fatal("failed to listen", zap.Error(err))
	}

	go func() {
		appLogger.Info("movie-service gRPC server started", zap.String("address", address))

		if err := grpcServer.Serve(listener); err != nil {
			appLogger.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	waitForShutdown(grpcServer, appLogger)
}

func waitForShutdown(grpcServer *grpc.Server, logger *zap.Logger) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	logger.Info("shutting down movie-service...")

	grpcServer.GracefulStop()

	logger.Info("movie-service stopped")
}
