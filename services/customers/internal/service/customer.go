package service

import (
	"context"
	"fmt"

	"github.com/Cladkoewka/marketplace-project/services/customers/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, c *domain.Customer) error
	GetByID(ctx context.Context, id int64) (*domain.Customer, error)
	GetByEmail(ctx context.Context, email string) (*domain.Customer, error)
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
