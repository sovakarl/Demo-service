package memory

import (
	"demo-service/internal/models"
)

type Cache struct {
	orders map[string]models.Order
}

func NewCache() *Cache {
	return &Cache{make(map[string]models.Order)}
}

func (c *Cache) Get(orderUID string) (*models.Order, error) {
	return nil, nil
}

func (c *Cache) Save(order *models.Order){
	//TODO ЖЕСТКОЕ
}
