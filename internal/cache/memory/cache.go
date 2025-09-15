package memory

import (
	"demo-service/internal/models"
	"hash/fnv"
	"log/slog"
	"time"
)

type memory struct {
	countShard        uint64
	chSignal          chan struct{}
	chStatusGC        chan struct{}
	shards            []shard
	defaultExpiration time.Duration
	defaultDuration   time.Duration
	logger            *slog.Logger
}

func NewCache(defaultExpiration, defaultDuration time.Duration, countShard uint64, logger *slog.Logger) *memory {
	if logger == nil {
		logger = slog.Default()
	}
	cache := memory{
		countShard:        countShard,
		chSignal:          make(chan struct{}),
		chStatusGC:        make(chan struct{},1),
		defaultExpiration: defaultExpiration,
		defaultDuration:   defaultDuration,
		shards:            make([]shard, countShard),
		logger:            logger.With("component", "cache"),
	}

	for i := range len(cache.shards) {
		cache.shards[i] = newShard()
	}

	if defaultDuration > 0 {
		cache.startGC()
	}
	logger.Info("cache start up ")
	return &cache
}

func (c *memory) hash(value string) uint64 {
	hasher := fnv.New64a()
	hasher.Write([]byte(value))
	index := hasher.Sum64() % uint64(c.countShard)
	return index
}

func (c *memory) Set(orderUID string, order *models.Order) {
	index := c.hash(orderUID)
	shard := c.getShard(index)
	c.logger.Debug("set order", "orderUID", orderUID, "shardIndex", index)
	shard.set(orderUID, order)
}

func (c *memory) Get(orderUID string) (*models.Order, bool) {
	index := c.hash(orderUID)
	shard := c.getShard(index)
	c.logger.Debug("get order", "orderUID", orderUID, "shardIndex", index)
	return shard.get(orderUID)
}

func (m *memory) getShard(index uint64) *shard {
	return &m.shards[index]
}

func (c *memory) Close() error {
	return c.stopGC()
}
