package service

import "demo-service/internal/models"

type Service interface {
	GetOrder(uid string) (*models.Order, error)
	SaveOrder(order *models.Order) error
}
