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
	rabbitmqadapter "github.com/luizdavid/movies-challenge/movie-service/internal/adapters/messaging/rabbitmq"
	mongoadapter "github.com/luizdavid/movies-challenge/movie-service/internal/adapters/repository/mongodb"
	"github.com/luizdavid/movies-challenge/movie-service/internal/bootstrap"
	"github.com/luizdavid/movies-challenge/movie-service/internal/config"
	"github.com/luizdavid/movies-challenge/movie-service/internal/database"
	"github.com/luizdavid/movies-challenge/movie-service/internal/logger"
	"github.com/luizdavid/movies-challenge/movie-service/internal/usecases"
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

	if err := bootstrap.EnsureMovieIndexes(ctx, collection, appLogger); err != nil {
		appLogger.Fatal("failed to ensure movie indexes", zap.Error(err))
	}

	if err := bootstrap.SeedMovies(ctx, collection, "movies.json", appLogger); err != nil {
		appLogger.Fatal("failed to seed movies", zap.Error(err))
	}

	rabbitConnection, err := rabbitmqadapter.NewConnection(cfg.RabbitMQURL)
	if err != nil {
		appLogger.Fatal("failed to connect to RabbitMQ", zap.Error(err))
	}
	defer rabbitConnection.Close()

	rabbitChannel, err := rabbitmqadapter.NewChannel(rabbitConnection)
	if err != nil {
		appLogger.Fatal("failed to create RabbitMQ channel", zap.Error(err))
	}
	defer rabbitChannel.Close()

	if err := rabbitmqadapter.DeclareMovieEventsExchange(rabbitChannel); err != nil {
		appLogger.Fatal("failed to declare RabbitMQ exchange", zap.Error(err))
	}

	movieEventPublisher := rabbitmqadapter.NewMovieEventPublisher(rabbitChannel)

	movieRepository := mongoadapter.NewMovieRepository(collection)
	movieUseCases := usecases.NewMovieUseCases(movieRepository, movieEventPublisher)
	movieHandler := grpcadapter.NewMovieHandler(movieUseCases)

	grpcServer := grpc.NewServer()

	moviepb.RegisterMovieServiceServer(grpcServer, movieHandler)

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
