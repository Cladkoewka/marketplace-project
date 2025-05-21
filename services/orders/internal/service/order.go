package service

import (
	"context"

	"github.com/Cladkoewka/marketplace-project/services/orders/internal/domain"
)

type OrderRepository interface {
	Create(context.Context, *domain.Order) error
	GetById(context.Context, int64) (*domain.Order, error)
	GetByCustomerId(context.Context, int64) ([]*domain.Order, error)
}

type OrderService struct {
	repository OrderRepository
}

func NewOrderService(repository OrderRepository) *OrderService {
	return &OrderService{repository: repository}
}

func (s *OrderService) Create(ctx context.Context, order *domain.Order) error {
	if order.Status == "" {
		order.Status = "pending"
	}

	return s.repository.Create(ctx, order)
}

func (s *OrderService) GetById(ctx context.Context, id int64) (*domain.Order, error) {
	return s.repository.GetById(ctx, id)
}

func (s *OrderService) GetByCustomerId(ctx context.Context, customerId int64) ([]*domain.Order, error) {
	return s.repository.GetByCustomerId(ctx, customerId)
}
