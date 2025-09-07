package repository

import "demo-service/internal/models"

type Repository interface {
	Get(orderUID string) (*models.Order, error)
	Insert(order *models.Order) error
}
