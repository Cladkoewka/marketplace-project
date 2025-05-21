package repository

import (
	"context"

	"github.com/Cladkoewka/marketplace-project/services/orders/internal/config"
	"github.com/Cladkoewka/marketplace-project/services/orders/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, cfg.DSN)
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	row := r.db.QueryRow(ctx, "INSERT INTO orders (customer_id, status) VALUES ($1, $2) RETURNING id", order.CustomerID, order.Status)
	return row.Scan(&order.ID)
}

func (r *OrderRepository) GetById(ctx context.Context, id int64) (*domain.Order, error) {
	order := &domain.Order{}
	row := r.db.QueryRow(ctx, "SELECT id, customer_id, status FROM orders WHERE id = $1", id)
	err := row.Scan(&order.ID, &order.CustomerID, &order.Status)
	return order, err
}

func (r *OrderRepository) GetByCustomerId(ctx context.Context, customerId int64) ([]*domain.Order, error) {
	var orders []*domain.Order

	rows, err := r.db.Query(ctx, "SELECT id, customer_id, status FROM orders WHERE customer_id = $1", customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := &domain.Order{}
		if err := rows.Scan(&order.ID, &order.CustomerID, &order.Status); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}