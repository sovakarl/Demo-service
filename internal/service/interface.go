package service

import (
	"context"
	"demo-service/internal/models"
)

type Service interface {
	GetOrder(ctx context.Context, uid string) (*models.Order, error)
	SaveOrder(ctx context.Context, order *models.Order) error
}
