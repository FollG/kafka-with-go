package models

import "time"

type EventType string

const (
	ProductCreated EventType = "product_created"
	ProductUpdated EventType = "product_updated"
	ProductDeleted EventType = "product_deleted"
)

type ProductEvent struct {
	EventID     string    `json:"event_id"`
	EventType   EventType `json:"event_type"`
	Timestamp   time.Time `json:"timestamp"`
	ProductID   int       `json:"product_id,omitempty"`
	ProductData *Product  `json:"product_data,omitempty"`

	ProducerID string `json:"producer_id"`
	Sequence   int64  `json:"sequence"`
}
