package cache

import "demo-service/internal/models"

type Cache interface {
	Save(order models.Order)
	Get(uid string) (models.Order, error)
}
