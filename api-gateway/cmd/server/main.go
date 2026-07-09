package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/adapters/grpcclient"
	httpadapter "github.com/luizdavid/movies-challenge/api-gateway/internal/adapters/http"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/config"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/logger"
	"github.com/luizdavid/movies-challenge/api-gateway/internal/routes"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
	router.Use(gin.Recovery())

	routes.RegisterRoutes(router, movieHandler)

	address := fmt.Sprintf(":%s", cfg.HTTPPort)

	appLogger.Info("api-gateway started", zap.String("address", address))

	if err := router.Run(address); err != nil {
		appLogger.Fatal("failed to start api-gateway", zap.Error(err))
	}
}
