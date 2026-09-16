package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"automation-test/orders/internal/domain"
	"automation-test/orders/internal/repository"
)

type OrderUsecase struct {
	repository repository.OrderRepository
}

func NewOrderUsecase(repository repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repository: repository}
}

func (u *OrderUsecase) Create(ctx context.Context, input domain.CreateOrderInput) (*domain.Order, error) {
	if err := validate(input); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	order := &domain.Order{
		CustomerName: strings.TrimSpace(input.CustomerName), Item: strings.TrimSpace(input.Item),
		Quantity: input.Quantity, TotalAmount: input.TotalAmount,
		Status: normalizeStatus(input.Status), CreatedAt: now, UpdatedAt: now,
	}
	if err := u.repository.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (u *OrderUsecase) List(ctx context.Context) ([]domain.Order, error) {
	return u.repository.List(ctx)
}

func (u *OrderUsecase) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	return u.repository.GetByID(ctx, id)
}

func (u *OrderUsecase) Update(ctx context.Context, id string, input domain.UpdateOrderInput) (*domain.Order, error) {
	if err := validate(input); err != nil {
		return nil, err
	}
	order, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	order.CustomerName = strings.TrimSpace(input.CustomerName)
	order.Item = strings.TrimSpace(input.Item)
	order.Quantity = input.Quantity
	order.TotalAmount = input.TotalAmount
	order.Status = normalizeStatus(input.Status)
	order.UpdatedAt = time.Now().UTC()
	if err := u.repository.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (u *OrderUsecase) Delete(ctx context.Context, id string) error {
	return u.repository.Delete(ctx, id)
}

func validate(input domain.CreateOrderInput) error {
	if strings.TrimSpace(input.CustomerName) == "" || strings.TrimSpace(input.Item) == "" {
		return fmt.Errorf("%w: customer_name and item are required", domain.ErrInvalidInput)
	}
	if input.Quantity <= 0 || input.TotalAmount < 0 {
		return fmt.Errorf("%w: quantity must be positive and total_amount cannot be negative", domain.ErrInvalidInput)
	}
	status := normalizeStatus(input.Status)
	if status != "pending" && status != "confirmed" && status != "cancelled" {
		return fmt.Errorf("%w: status must be pending, confirmed, or cancelled", domain.ErrInvalidInput)
	}
	return nil
}

func normalizeStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return "pending"
	}
	return strings.ToLower(strings.TrimSpace(status))
}
