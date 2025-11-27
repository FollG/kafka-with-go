package services

import (
	"fmt"
	"time"

	"github.com/FollG/kafka-with-go/internal/domain/models"
)

type ProductValidator struct{}

func NewProductValidator() *ProductValidator {
	return &ProductValidator{}
}

func (v *ProductValidator) Validate(product *models.Product) error {
	// Базовые проверки
	if product.Name == "" || len(product.Name) > 255 {
		return fmt.Errorf("product name must be between 1 and 255 characters")
	}

	if product.Weight <= 0 {
		return fmt.Errorf("product weight must be positive")
	}

	if !isValidUnit(product.Unit) {
		return fmt.Errorf("invalid unit: %s. Must be one of: g, kg, l, piece", product.Unit)
	}

	if product.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}

	if len(product.Color) > 50 {
		return fmt.Errorf("color must be 50 characters or less")
	}

	// Валидация по типу продукта
	switch product.Type {
	case models.ClothingHeadwear:
		return v.validateHeadwear(product)
	case models.ClothingBody:
		return v.validateBodyClothing(product)
	case models.ClothingPants:
		return v.validatePants(product)
	case models.ClothingShoes:
		return v.validateShoes(product)
	case models.Electronics:
		return v.validateElectronics(product)
	case models.Food:
		return v.validateFood(product)
	case models.Furniture, models.HomeGoods:
		return v.validateFurniture(product)
	case models.Adult:
		return nil // Для товаров 18+ особых валидаций нет
	default:
		return fmt.Errorf("unknown product type: %s", product.Type)
	}
}

func (v *ProductValidator) validateHeadwear(product *models.Product) error {
	if product.Attributes.HeadCircumference == nil {
		return fmt.Errorf("head circumference is required for headwear")
	}
	if *product.Attributes.HeadCircumference <= 0 {
		return fmt.Errorf("head circumference must be positive")
	}
	if product.Unit != "piece" {
		return fmt.Errorf("headwear must be measured in pieces")
	}
	return nil
}

func (v *ProductValidator) validateBodyClothing(product *models.Product) error {
	if product.Attributes.ChestCircumference == nil {
		return fmt.Errorf("chest circumference is required for body clothing")
	}
	if *product.Attributes.ChestCircumference <= 0 {
		return fmt.Errorf("chest circumference must be positive")
	}
	if product.Unit != "piece" {
		return fmt.Errorf("body clothing must be measured in pieces")
	}
	return nil
}

func (v *ProductValidator) validatePants(product *models.Product) error {
	if product.Attributes.WaistCircumference == nil {
		return fmt.Errorf("waist circumference is required for pants")
	}
	if *product.Attributes.WaistCircumference <= 0 {
		return fmt.Errorf("waist circumference must be positive")
	}
	if product.Unit != "piece" {
		return fmt.Errorf("pants must be measured in pieces")
	}
	return nil
}

func (v *ProductValidator) validateShoes(product *models.Product) error {
	if product.Attributes.FootSize == nil {
		return fmt.Errorf("foot size is required for shoes")
	}
	if *product.Attributes.FootSize <= 0 {
		return fmt.Errorf("foot size must be positive")
	}
	if product.Unit != "piece" {
		return fmt.Errorf("shoes must be measured in pieces")
	}
	return nil
}

func (v *ProductValidator) validateElectronics(product *models.Product) error {
	if product.Attributes.WarrantyMonths == nil {
		return fmt.Errorf("warranty months is required for electronics")
	}
	if *product.Attributes.WarrantyMonths <= 0 {
		return fmt.Errorf("warranty months must be positive")
	}
	return nil
}

func (v *ProductValidator) validateFood(product *models.Product) error {
	if product.Attributes.ExpiryDate == nil {
		return fmt.Errorf("expiry date is required for food")
	}
	if product.Attributes.ExpiryDate.Before(time.Now()) {
		return fmt.Errorf("expiry date cannot be in the past")
	}
	return nil
}

func (v *ProductValidator) validateFurniture(product *models.Product) error {
	// Для мебели проверяем, что указаны размеры
	if product.Attributes.Dimensions == "" {
		return fmt.Errorf("dimensions are required for furniture")
	}
	return nil
}

func isValidUnit(unit string) bool {
	validUnits := map[string]bool{
		"g":     true,
		"kg":    true,
		"l":     true,
		"piece": true,
	}
	return validUnits[unit]
}
