package repository

import (
	"context"

	"github.com/Cladkoewka/marketplace-project/services/customers/internal/config"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, cfg.DSN)
}

func NewCustomerRepository(db *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(ctx context.Context, c *domain.Customer) error {
	row := r.db.QueryRow(ctx, `
    INSERT INTO customers (email, name, phone, address) 
    VALUES ($1, $2, $3, $4) RETURNING id
	`, c.Email, c.Name, c.Phone, c.Address)
	return row.Scan(&c.ID)
}

func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*domain.Customer, error) {
	c := &domain.Customer{}
	row := r.db.QueryRow(ctx, `
    SELECT id, email, name, phone, address FROM customers WHERE id = $1
	`, id)
	err := row.Scan(&c.ID, &c.Email, &c.Name, &c.Phone, &c.Address)
	return c, err
}

func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	c := &domain.Customer{}
	row := r.db.QueryRow(ctx, `SELECT id, email, name, phone, address FROM customers WHERE email = $1`, email)
	err := row.Scan(&c.ID, &c.Email, &c.Name, &c.Phone, &c.Address)
	return c, err
}
