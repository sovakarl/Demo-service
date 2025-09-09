package repository

import (
	"demo-service/internal/models"
	"io"
)

type Repository interface {
	Get(orderUID string) (*models.Order, error)
	GetAll(rowsCount uint64) ([]*models.Order, error)
	Insert(order *models.Order) error
	io.Closer
}
