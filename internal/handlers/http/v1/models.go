package v1

import (
	"time"
)

type CreateProductRequest struct {
	Name       string            `json:"name" validate:"required,min=1,max=255"`
	Weight     float64           `json:"weight" validate:"required,gt=0"`
	Unit       string            `json:"unit" validate:"required,oneof=g kg l piece"`
	Color      string            `json:"color" validate:"max=50"`
	Type       string            `json:"type" validate:"required,oneof=clothing_headwear clothing_body clothing_pants clothing_shoes food furniture electronics adult home_goods"`
	Price      float64           `json:"price" validate:"required,gt=0"`
	Attributes AttributesRequest `json:"attributes"`
}

type UpdateProductRequest struct {
	Name       string            `json:"name" validate:"required,min=1,max=255"`
	Weight     float64           `json:"weight" validate:"required,gt=0"`
	Unit       string            `json:"unit" validate:"required,oneof=g kg l piece"`
	Color      string            `json:"color" validate:"max=50"`
	Type       string            `json:"type" validate:"required,oneof=clothing_headwear clothing_body clothing_pants clothing_shoes food furniture electronics adult home_goods"`
	Price      float64           `json:"price" validate:"required,gt=0"`
	Attributes AttributesRequest `json:"attributes"`
}

type AttributesRequest struct {
	Size               string     `json:"size,omitempty"`
	HeadCircumference  *float64   `json:"head_circumference,omitempty"`
	ChestCircumference *float64   `json:"chest_circumference,omitempty"`
	WaistCircumference *float64   `json:"waist_circumference,omitempty"`
	HipCircumference   *float64   `json:"hip_circumference,omitempty"`
	FootSize           *float64   `json:"foot_size,omitempty"`
	ExpiryDate         *time.Time `json:"expiry_date,omitempty"`
	NutritionalInfo    string     `json:"nutritional_info,omitempty"`
	WarrantyMonths     *int       `json:"warranty_months,omitempty"`
	Voltage            string     `json:"voltage,omitempty"`
	Dimensions         string     `json:"dimensions,omitempty"`
	Material           string     `json:"material,omitempty"`
}

type ProductResponse struct {
	ID         int                `json:"id"`
	Name       string             `json:"name"`
	Weight     float64            `json:"weight"`
	Unit       string             `json:"unit"`
	Color      string             `json:"color"`
	Type       string             `json:"type"`
	Price      float64            `json:"price"`
	Attributes AttributesResponse `json:"attributes"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type AttributesResponse struct {
	Size               string     `json:"size,omitempty"`
	HeadCircumference  *float64   `json:"head_circumference,omitempty"`
	ChestCircumference *float64   `json:"chest_circumference,omitempty"`
	WaistCircumference *float64   `json:"waist_circumference,omitempty"`
	HipCircumference   *float64   `json:"hip_circumference,omitempty"`
	FootSize           *float64   `json:"foot_size,omitempty"`
	ExpiryDate         *time.Time `json:"expiry_date,omitempty"`
	NutritionalInfo    string     `json:"nutritional_info,omitempty"`
	WarrantyMonths     *int       `json:"warranty_months,omitempty"`
	Voltage            string     `json:"voltage,omitempty"`
	Dimensions         string     `json:"dimensions,omitempty"`
	Material           string     `json:"material,omitempty"`
}

type CreateProductResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type UpdateProductResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DeleteProductResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ListProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
