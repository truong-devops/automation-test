package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("payment not found")
	ErrInvalidID    = errors.New("invalid payment id")
	ErrInvalidInput = errors.New("invalid payment input")
)

type Payment struct {
	ID        string    `json:"id" bson:"-"`
	OrderID   string    `json:"order_id" bson:"order_id"`
	Amount    float64   `json:"amount" bson:"amount"`
	Method    string    `json:"method" bson:"method"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CreatePaymentInput struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
	Status  string  `json:"status"`
}

type UpdatePaymentInput = CreatePaymentInput
