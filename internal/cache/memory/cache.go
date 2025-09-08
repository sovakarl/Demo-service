package memory

import (
	"demo-service/internal/models"
	"hash/fnv"
	"time"
)

type memory struct {
	countShard        uint
	chSignal          chan struct{}
	shards            []shard
	defaultExpiration time.Duration
	defaultDuration   time.Duration
}

func NewCache(defaultExpiration, defaultDuration time.Duration, countShard uint) *memory {
	cache := memory{
		countShard:        countShard,
		chSignal:          make(chan struct{}),
		defaultExpiration: defaultExpiration,
		defaultDuration:   defaultDuration,
		shards:            make([]shard,countShard),
	}

	for i := range len(cache.shards) {
		cache.shards[i] = newShard()
	}

	if defaultDuration > 0 {
		cache.startGC()
	}
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
	shard.set(orderUID, order)
}

func (c *memory) Get(orderUID string) (*models.Order, bool) {
	index := c.hash(orderUID)
	shard := c.getShard(index)
	return shard.get(orderUID)
}

func (m *memory) getShard(index uint64) *shard {
	return &m.shards[index]
}

func (c *memory) Close() {
	c.stopGC()
}
