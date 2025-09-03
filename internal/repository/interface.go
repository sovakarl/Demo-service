package repository

import "demo-service/internal/models"

type Repository interface {
	Get(uid string) (*models.Order, error)
	Insert(order *models.Order) error
}
