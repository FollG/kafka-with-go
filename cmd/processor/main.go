package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FollG/kafka-with-go/internal/adapters/kafka"
	"github.com/FollG/kafka-with-go/internal/adapters/postgres"
	"github.com/FollG/kafka-with-go/internal/adapters/redis"
	"github.com/FollG/kafka-with-go/internal/pkg/cache"
	"github.com/FollG/kafka-with-go/internal/pkg/config"
	"github.com/FollG/kafka-with-go/internal/pkg/database"
	"github.com/FollG/kafka-with-go/internal/pkg/logger"
	"github.com/FollG/kafka-with-go/internal/pkg/metrics"
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
			panic(err)
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
			panic(err)
		}
	}(redisClient)

	// repo init
	productRepo := postgres.NewProductRepository(db)
	productCache := redis.NewProductCache(redisClient, cfg.Redis.TTL)

	// kafka consumer
	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.Topic,
		cfg.Kafka.ConsumerGroup,
		100,                  // batchSize
		100*time.Millisecond, // batchTimeout
	)
	consumer.SetProductRepo(productRepo)
	consumer.SetCache(productCache)
	defer func(consumer *kafka.Consumer) {
		err := consumer.Close()
		if err != nil {
			panic(err)
		}
	}(consumer)

	// graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Info(ctx, "starting Kafka consumer",
			"topic", cfg.Kafka.Topic,
			"group", cfg.Kafka.ConsumerGroup,
		)
		if err := consumer.Start(ctx); err != nil {
			logger.Fatal(ctx, "failed to start consumer", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "shutting down processor...")

	cancel()

	time.Sleep(5 * time.Second)

	logger.Info(ctx, "processor exited")
}
