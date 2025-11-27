package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/FollG/kafka-with-go/internal/domain/models"
	"github.com/FollG/kafka-with-go/internal/domain/repositories"
	vld "github.com/FollG/kafka-with-go/internal/pkg/validator"
)

type ProductUseCase struct {
	repo          repositories.ProductRepository
	cache         repositories.ProductCache
	eventProducer repositories.EventProducer
	validator     *vld.ProductValidator
}

func NewProductUseCase(
	repo repositories.ProductRepository,
	cache repositories.ProductCache,
	eventProducer repositories.EventProducer,
	validator *vld.ProductValidator,
) *ProductUseCase {
	return &ProductUseCase{
		repo:          repo,
		cache:         cache,
		eventProducer: eventProducer,
		validator:     validator,
	}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, product *models.Product) error {
	if err := uc.validator.Validate(product); err != nil {
		return err
	}

	event := &models.ProductEvent{
		EventID:     generateEventID(),
		EventType:   models.ProductCreated,
		Timestamp:   time.Now(),
		ProductData: product,
		ProducerID:  "product-api",
		Sequence:    time.Now().UnixNano(),
	}

	if err := uc.eventProducer.SendProductEvent(ctx, event); err != nil {
		return err
	}

	return nil
}

func (uc *ProductUseCase) GetProduct(ctx context.Context, id int) (*models.Product, error) {
	// Пытаемся получить из кеша
	cacheKey := fmt.Sprintf("product:%d", id)
	if cached, err := uc.cache.Get(ctx, cacheKey); err == nil && cached != nil {
		return cached, nil
	}

	// Если нет в кеше, идем в базу
	product, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// На случай, если реализация репозитория вернет (nil, nil)
	if product == nil {
		return nil, models.ErrProductNotFound
	}

	// Сохраняем в кеш
	if err := uc.cache.Set(ctx, cacheKey, product); err != nil {
		// Логируем ошибку, но не прерываем выполнение
		fmt.Printf("Failed to cache product: %v\n", err)
	}

	return product, nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *models.Product) error {
	// Валидация
	if err := uc.validator.Validate(product); err != nil {
		return err
	}

	// Создаем событие для Kafka
	event := &models.ProductEvent{
		EventID:     generateEventID(),
		EventType:   models.ProductUpdated,
		Timestamp:   time.Now(),
		ProductID:   product.ID,
		ProductData: product,
		ProducerID:  "product-api",
		Sequence:    time.Now().UnixNano(),
	}

	if err := uc.eventProducer.SendProductEvent(ctx, event); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("product:%d", product.ID)
	if err := uc.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("Failed to invalidate cache: %v\n", err)
	}

	return nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id int) error {
	event := &models.ProductEvent{
		EventID:    generateEventID(),
		EventType:  models.ProductDeleted,
		Timestamp:  time.Now(),
		ProductID:  id,
		ProducerID: "product-api",
		Sequence:   time.Now().UnixNano(),
	}

	if err := uc.eventProducer.SendProductEvent(ctx, event); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	if err := uc.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("Failed to invalidate cache: %v\n", err)
	}

	return nil
}

func (uc *ProductUseCase) ListProducts(ctx context.Context, filter models.ProductFilter) ([]*models.Product, error) {
	// для списков также можно использовать кеш, но это сложнее из-за вариативности фильтров
	return uc.repo.List(ctx, filter)
}

func generateEventID() string {
	return fmt.Sprintf("event-%d", time.Now().UnixNano())
}
