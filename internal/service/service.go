package service

import (
	"demo-service/internal/cache"
	"demo-service/internal/models"
	"demo-service/internal/repository"
	"log/slog"
)

type OrderService struct {
	db    repository.Repository
	cache cache.Cache
	cnf   Config
	log   *slog.Logger
}

type Config struct {
	CacheWarmUpLimit uint64
}

func (s *OrderService) updateCache() {
	defer func() {
		if r := recover(); r != nil {
			s.log.Error("паника при прогреве кеша", r)
		}
	}()
	orders, err := s.db.GetAll(s.cnf.CacheWarmUpLimit)
	if err != nil {
		s.log.Error("не удалось загрузить заказы для прогрева кеша", "error", err)
		return
	}
	for _, order := range orders {
		orderUID := order.GetUid()
		s.cache.Set(orderUID, order)
	}
}

func NewService(db repository.Repository, cache cache.Cache, cnf Config, log *slog.Logger) Service {
	if log == nil {
		log = slog.Default()
	}
	if cnf.CacheWarmUpLimit == 0 {
		cnf.CacheWarmUpLimit = 100
	}

	service := &OrderService{
		db:    db,
		cache: cache,
		cnf:   cnf,
		log:   log.With("component", "order_service"),
	}
	go service.updateCache()
	return service
}

func (s *OrderService) GetOrder(uid string) (*models.Order, error) {
	if order, ex := s.cache.Get(uid); ex {
		return order, nil
	}
	order, err := s.db.Get(uid)
	if err != nil {
		return nil, err
	}
	s.cache.Set(uid, order)
	return order, nil
}

func (s *OrderService) SaveOrder(order *models.Order) error {
	err := s.db.Insert(order)
	return err
}
