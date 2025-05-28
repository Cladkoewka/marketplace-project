package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Cladkoewka/marketplace-project/services/customers/internal/domain"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/kafka/event"
)

type CustomerRepository interface {
	Create(ctx context.Context, c *domain.Customer) error
	GetByID(ctx context.Context, id int64) (*domain.Customer, error)
	GetByEmail(ctx context.Context, email string) (*domain.Customer, error)
	IncrementActivityScore(ctx context.Context, customerID int64, delta int) error
}

type CustomerService struct {
	repo CustomerRepository
}

func NewCustomerService(repo CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(ctx context.Context, c *domain.Customer) error {
	if c.Email == "" {
		return fmt.Errorf("email is empty")
	}
	
	return s.repo.Create(ctx, c)
}

func (s *CustomerService) GetByID(ctx context.Context, id int64) (*domain.Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CustomerService) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *CustomerService) HandleOrderPlaced(ctx context.Context, e event.OrderPlacedEvent) error {
	slog.Info("Received OrderPlacedEvent", "event", e)

	err := s.repo.IncrementActivityScore(ctx, e.CustomerID, 10)
	if err != nil {
		return fmt.Errorf("incrementing activity score: %w", err)
	}

	slog.Info("Customer activity score incremented",
		"customer_id", e.CustomerID,
		"delta", 10,
	)
	return nil
}