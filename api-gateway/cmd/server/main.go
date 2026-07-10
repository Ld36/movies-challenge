package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/luizdavid/movies-challenge/api-gateway/docs"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/adapters/grpcclient"
	httpadapter "github.com/luizdavid/movies-challenge/api-gateway/internal/adapters/http"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/config"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/logger"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/middleware"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/routes"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// @title API Gateway de Filmes
// @version 1.0
// @description API Gateway para o microsserviço de Filmes utilizando HTTP, gRPC, Go e MongoDB.
// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.Load()

	appLogger, err := logger.New(cfg.AppEnv)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	conn, err := grpc.NewClient(
		cfg.MovieServiceGRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		appLogger.Fatal("failed to connect to movie-service", zap.Error(err))
	}
	defer conn.Close()

	movieClient := grpcclient.NewMovieClient(conn)
	movieHandler := httpadapter.NewMovieHandler(movieClient)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		middleware.RequestID(),
		middleware.Logger(appLogger),
		gin.Recovery(),
	)

	routes.RegisterRoutes(router, movieHandler)

	address := fmt.Sprintf(":%s", cfg.HTTPPort)

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	go func() {
		appLogger.Info("api-gateway started", zap.String("address", address))

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("failed to start api-gateway", zap.Error(err))
		}
	}()

	waitForShutdown(server, appLogger)
}

func waitForShutdown(server *http.Server, appLogger *zap.Logger) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	appLogger.Info("shutting down api-gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("api-gateway forced to shutdown", zap.Error(err))
		return
	}

	appLogger.Info("api-gateway stopped gracefully")
}
