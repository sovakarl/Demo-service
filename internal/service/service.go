package service

import (
	"demo-service/internal/cache"
	"demo-service/internal/models"
	"demo-service/internal/repository"
)

type OrderService struct {
	db    repository.Repository
	cache cache.Cache
}

func NewService(db repository.Repository, cache cache.Cache) Service {
	return OrderService{db: db, cache: cache}
}

func (s OrderService) GetOrder(uid string) (*models.Order, error) {
	if order, ex := s.cache.Get(uid); ex {
		return order, nil
	}
	order, err := s.db.Get(uid)
	if err != nil {
		return nil, err
	}
	s.cache.Set(order)
	return order, nil
}

func (s OrderService) SaveOrder(order *models.Order) error {
	err := s.db.Insert(order)
	return err
}
