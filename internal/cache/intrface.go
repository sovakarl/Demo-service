package cache

import "demo-service/internal/models"

type Cache interface {
	Set(order *models.Order)
	Get(orderUID string) (*models.Order, bool)
	Delete(orderUID string)
}
