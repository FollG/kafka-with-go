package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FollG/kafka-with-go/internal/domain/models"

	"github.com/redis/go-redis/v9"
)

type ProductCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewProductCache(client *redis.Client, ttl time.Duration) *ProductCache {
	return &ProductCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *ProductCache) Get(ctx context.Context, key string) (*models.Product, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Ключ не найден - это не ошибка
		}
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var product models.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %w", err)
	}

	return &product, nil
}

func (c *ProductCache) Set(ctx context.Context, key string, product *models.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	err = c.client.Set(ctx, key, data, c.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (c *ProductCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

func (c *ProductCache) SetList(ctx context.Context, key string, products []*models.Product) error {
	data, err := json.Marshal(products)
	if err != nil {
		return fmt.Errorf("failed to marshal products list: %w", err)
	}

	err = c.client.Set(ctx, key, data, c.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set list cache: %w", err)
	}

	return nil
}

func (c *ProductCache) GetList(ctx context.Context, key string) ([]*models.Product, error) {
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Ключ не найден - это не ошибка
		}
		return nil, fmt.Errorf("failed to get list from cache: %w", err)
	}

	var products []*models.Product
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, fmt.Errorf("failed to unmarshal products list: %w", err)
	}

	return products, nil
}
