package service

import (
	"context"
	"demo-service/internal/cache"
	"demo-service/internal/models"
	"demo-service/internal/repository"
	"log/slog"
	"runtime/debug"
	"time"
)

type OrderService struct {
	db     repository.Repository
	cache  cache.Cache
	cnf    Config
	logger *slog.Logger
}

type Config struct {
	CacheWarmUpLimit uint64
}

func NewService(db repository.Repository, cache cache.Cache, cnf Config, log *slog.Logger) Service {
	if log == nil {
		log = slog.Default()
	}
	if cnf.CacheWarmUpLimit == 0 {
		cnf.CacheWarmUpLimit = 100
	}

	service := &OrderService{
		db:     db,
		cache:  cache,
		cnf:    cnf,
		logger: log.With("component", "order_service"),
	}
	go service.updateCache()
	return service
}

func (s *OrderService) updateCache() {
	s.logger.Info("start cache warm-up ")
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic when warm-up the cache", r)
			debug.PrintStack()
		}
	}()

	ctx, done := context.WithTimeout(context.Background(), time.Minute)
	defer done()

	orders, err := s.db.GetAll(ctx, s.cnf.CacheWarmUpLimit)
	if err != nil {
		s.logger.Debug("Failed to load orders for cache warm-up", "error", err)
		return
	}
	if len(orders) == 0 {
		s.logger.Debug("no orders found for cache warm-up")
		return
	}
	for _, order := range orders {
		if orderUID := order.GetUid(); orderUID != "" {
			s.cache.Set(orderUID, order)
		}
	}

	s.logger.Info("fetched orders for cache warm-up",
		"count", len(orders),
	)
}

func (s *OrderService) GetOrder(ctx context.Context, uid string) (*models.Order, error) {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Debug("GetOrder completed ",
			"order_uid", uid,
			"duration_ms", duration.Milliseconds(),
		)
	}(time.Now())

	if order, ex := s.cache.Get(uid); ex {
		return order, nil
	}

	ctx, done := context.WithTimeout(ctx, time.Minute)
	defer done()

	order, err := s.db.Get(ctx, uid)
	if err != nil {
		return nil, err
	}
	s.cache.Set(uid, order)
	return order, nil
}

func (s *OrderService) SaveOrder(ctx context.Context, order *models.Order) error {
	defer func(start time.Time) {
		duration := time.Since(start)
		s.logger.Debug("GetOrder completed ",
			"order_uid", order.GetUid(),
			"duration_ms", duration.Milliseconds(),
		)
	}(time.Now())
	
	ctx, done := context.WithTimeout(ctx, time.Minute)
	defer done()

	err := s.db.Insert(ctx, order)
	return err
}
