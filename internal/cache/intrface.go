package cache

import "demo-service/internal/models"

type Cache interface {
	Set(orderUID string, order  *models.Order)
	Get(orderUID string) (*models.Order, bool)
}
