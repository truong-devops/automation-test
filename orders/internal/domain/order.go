package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("order not found")
	ErrInvalidID    = errors.New("invalid order id")
	ErrInvalidInput = errors.New("invalid order input")
)

type Order struct {
	ID           string    `json:"id" bson:"-"`
	CustomerName string    `json:"customer_name" bson:"customer_name"`
	Item         string    `json:"item" bson:"item"`
	Quantity     int       `json:"quantity" bson:"quantity"`
	TotalAmount  float64   `json:"total_amount" bson:"total_amount"`
	Status       string    `json:"status" bson:"status"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateOrderInput struct {
	CustomerName string  `json:"customer_name"`
	Item         string  `json:"item"`
	Quantity     int     `json:"quantity"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status"`
}

type UpdateOrderInput = CreateOrderInput
