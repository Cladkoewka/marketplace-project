package service

import (
	"context"
	"time"

	"github.com/Cladkoewka/marketplace-project/services/orders/internal/domain"
	"github.com/Cladkoewka/marketplace-project/services/orders/internal/kafka"
	"github.com/Cladkoewka/marketplace-project/services/orders/internal/kafka/event"
)

type OrderRepository interface {
	Create(context.Context, *domain.Order) error
	GetById(context.Context, int64) (*domain.Order, error)
	GetByCustomerId(context.Context, int64) ([]*domain.Order, error)
}

type OrderService struct {
	repository OrderRepository
	producer   *kafka.Producer
}

func NewOrderService(repo OrderRepository, producer *kafka.Producer) *OrderService {
	return &OrderService{repository: repo, producer: producer}
}

func (s *OrderService) Create(ctx context.Context, order *domain.Order) error {
	if order.Status == "" {
		order.Status = "pending"
	}

	if err := s.repository.Create(ctx, order); err != nil {
		return err
	}

	event := event.OrderPlacedEvent{
		Event:      "order_placed",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
	}

	return s.producer.Send(ctx, "order", event)
}

func (s *OrderService) GetById(ctx context.Context, id int64) (*domain.Order, error) {
	return s.repository.GetById(ctx, id)
}

func (s *OrderService) GetByCustomerId(ctx context.Context, customerId int64) ([]*domain.Order, error) {
	return s.repository.GetByCustomerId(ctx, customerId)
}
