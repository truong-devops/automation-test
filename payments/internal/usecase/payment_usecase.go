package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"automation-test/payments/internal/domain"
	"automation-test/payments/internal/repository"
)

type PaymentUsecase struct {
	repository repository.PaymentRepository
}

func NewPaymentUsecase(repository repository.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{repository: repository}
}

func (u *PaymentUsecase) Create(ctx context.Context, input domain.CreatePaymentInput) (*domain.Payment, error) {
	if err := validate(input); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	payment := &domain.Payment{
		OrderID: strings.TrimSpace(input.OrderID), Amount: input.Amount,
		Method: strings.ToLower(strings.TrimSpace(input.Method)), Status: normalizeStatus(input.Status),
		CreatedAt: now, UpdatedAt: now,
	}
	if err := u.repository.Create(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (u *PaymentUsecase) List(ctx context.Context) ([]domain.Payment, error) {
	return u.repository.List(ctx)
}

func (u *PaymentUsecase) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	return u.repository.GetByID(ctx, id)
}

func (u *PaymentUsecase) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	return u.repository.GetByOrderID(ctx, orderID)
}

func (u *PaymentUsecase) Update(ctx context.Context, id string, input domain.UpdatePaymentInput) (*domain.Payment, error) {
	if err := validate(input); err != nil {
		return nil, err
	}
	payment, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	payment.OrderID = strings.TrimSpace(input.OrderID)
	payment.Amount = input.Amount
	payment.Method = strings.ToLower(strings.TrimSpace(input.Method))
	payment.Status = normalizeStatus(input.Status)
	payment.UpdatedAt = time.Now().UTC()
	if err := u.repository.Update(ctx, payment); err != nil {
		return nil, err
	}
	return payment, nil
}

func (u *PaymentUsecase) Delete(ctx context.Context, id string) error {
	return u.repository.Delete(ctx, id)
}

func validate(input domain.CreatePaymentInput) error {
	if strings.TrimSpace(input.OrderID) == "" || strings.TrimSpace(input.Method) == "" {
		return fmt.Errorf("%w: order_id and method are required", domain.ErrInvalidInput)
	}
	if input.Amount <= 0 {
		return fmt.Errorf("%w: amount must be positive", domain.ErrInvalidInput)
	}
	method := strings.ToLower(strings.TrimSpace(input.Method))
	if method != "card" && method != "bank_transfer" && method != "cash" {
		return fmt.Errorf("%w: method must be card, bank_transfer, or cash", domain.ErrInvalidInput)
	}
	status := normalizeStatus(input.Status)
	if status != "pending" && status != "completed" && status != "failed" && status != "refunded" {
		return fmt.Errorf("%w: status must be pending, completed, failed, or refunded", domain.ErrInvalidInput)
	}
	return nil
}

func normalizeStatus(status string) string {
	if strings.TrimSpace(status) == "" {
		return "pending"
	}
	return strings.ToLower(strings.TrimSpace(status))
}
