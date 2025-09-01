package repository

import "demo-service/internal/models"

type Repository interface {
	Get(uid string) (*models.Order, error)
	Put(order *models.Order) error
}
