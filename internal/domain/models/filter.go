package models

type ProductFilter struct {
	Limit    int
	Offset   int
	MinPrice *float64
	MaxPrice *float64
	Color    string
	Types    []ProductType
}
