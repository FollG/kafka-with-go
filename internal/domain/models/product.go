package models

import (
	"encoding/json"
	"time"
)

type ProductType string

const (
	ClothingHeadwear ProductType = "clothing_headwear"
	ClothingBody     ProductType = "clothing_body"
	ClothingPants    ProductType = "clothing_pants"
	ClothingShoes    ProductType = "clothing_shoes"
	Food             ProductType = "food"
	Furniture        ProductType = "furniture"
	Electronics      ProductType = "electronics"
	Adult            ProductType = "adult"
	HomeGoods        ProductType = "home_goods"
)

type Product struct {
	ID         int         `json:"id"`
	Name       string      `json:"name"`
	Weight     float64     `json:"weight"`
	Unit       string      `json:"unit"` // "g", "kg", "l", "piece"
	Color      string      `json:"color"`
	Type       ProductType `json:"type"`
	Price      float64     `json:"price"`
	Attributes Attributes  `json:"attributes"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type Attributes struct {
	Size string `json:"size,omitempty"` // S, M, L, XL, 42, 44, etc.

	HeadCircumference *float64 `json:"head_circumference,omitempty"`

	ChestCircumference *float64 `json:"chest_circumference,omitempty"`

	WaistCircumference *float64 `json:"waist_circumference,omitempty"`
	HipCircumference   *float64 `json:"hip_circumference,omitempty"`

	FootSize *float64 `json:"foot_size,omitempty"`

	ExpiryDate      *time.Time `json:"expiry_date,omitempty"`
	NutritionalInfo string     `json:"nutritional_info,omitempty"`

	WarrantyMonths *int   `json:"warranty_months,omitempty"`
	Voltage        string `json:"voltage,omitempty"`

	Dimensions string `json:"dimensions,omitempty"` // "100x50x200"
	Material   string `json:"material,omitempty"`
}

func (a *Attributes) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), a)
}

func (a Attributes) Value() ([]byte, error) {
	return json.Marshal(a)
}
