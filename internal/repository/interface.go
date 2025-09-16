package repository

import (
	"context"
	"demo-service/internal/models"
)

type Repository interface {
	Get(ctx context.Context, orderUID string) (*models.Order, error)
	GetAll(ctx context.Context, rowsCount uint64) ([]*models.Order, error)
	Insert(ctx context.Context, order *models.Order) error
}
