package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FollG/kafka-with-go/internal/adapters/kafka"
	"github.com/FollG/kafka-with-go/internal/adapters/postgres"
	"github.com/FollG/kafka-with-go/internal/adapters/redis"
	"github.com/FollG/kafka-with-go/internal/domain/services"
	"github.com/FollG/kafka-with-go/internal/handlers/http/v1"
	"github.com/FollG/kafka-with-go/internal/pkg/cache"
	"github.com/FollG/kafka-with-go/internal/pkg/config"
	"github.com/FollG/kafka-with-go/internal/pkg/database"
	"github.com/FollG/kafka-with-go/internal/pkg/logger"
	"github.com/FollG/kafka-with-go/internal/pkg/metrics"
	vld "github.com/FollG/kafka-with-go/internal/pkg/validator"
	"github.com/FollG/kafka-with-go/internal/usecases"
	redis2 "github.com/redis/go-redis/v9"
)

func main() {
	// conf
	cfg := config.Load()

	// logger
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// metrics
	metrics.Init(cfg.Metrics.Port)

	// psql
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		logger.Fatal(context.Background(), "failed to connect to database", "error", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(context.Background(), "failed to close database connection", "error", err)
		}
	}(db)

	// redis
	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Fatal(context.Background(), "failed to connect to redis", "error", err)
	}
	defer func(redisClient *redis2.Client) {
		err := redisClient.Close()
		if err != nil {
			logger.Error(context.Background(), "failed to close redis client", "error", err)
		}
	}(redisClient)

	// kafka producer
	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer func(producer *kafka.Producer) {
		err := producer.Close()
		if err != nil {
			logger.Error(context.Background(), "failed to close producer", "error", err)
		}
	}(producer)

	// reps and services
	productRepo := postgres.NewProductRepository(db)
	productCache := redis.NewProductCache(redisClient, cfg.Redis.TTL)
	validator := services.NewProductValidator()

	// usecases
	productUC := usecases.NewProductUseCase(productRepo, productCache, producer, (*vld.ProductValidator)(validator))

	// http server
	router := v1.NewRouter(productUC, db, redisClient, cfg.Server.RateLimit)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info(context.Background(), "starting API server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(context.Background(), "failed to start server", "error", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(context.Background(), "shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error(context.Background(), "server forced to shutdown", "error", err)
	}

	logger.Info(context.Background(), "server exited")
}
