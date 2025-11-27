package repositories

import (
	"context"

	"github.com/FollG/kafka-with-go/internal/domain/models"
)

// ProductRepository определяет контракт для работы с продуктами в БД
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id int) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter models.ProductFilter) ([]*models.Product, error)
}

// ProductCache определяет контракт для кеширования продуктов
type ProductCache interface {
	Get(ctx context.Context, key string) (*models.Product, error)
	Set(ctx context.Context, key string, product *models.Product) error
	Delete(ctx context.Context, key string) error
	SetList(ctx context.Context, key string, products []*models.Product) error
	GetList(ctx context.Context, key string) ([]*models.Product, error)
}

// EventProducer определяет контракт для отправки событий в Kafka
type EventProducer interface {
	SendProductEvent(ctx context.Context, event *models.ProductEvent) error
	Close() error
}

// EventConsumer определяет контракт для потребления событий из Kafka
type EventConsumer interface {
	Start(ctx context.Context) error
	Close() error
}
