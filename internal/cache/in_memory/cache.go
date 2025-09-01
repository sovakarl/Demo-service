package memory

import "demo-service/internal/models"

type Cache struct {
	orders map[string]models.Order
}

func NewCache() Cache {
	return Cache{make(map[string]models.Order)}
}
