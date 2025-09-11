package memory

import (
	"demo-service/internal/models"
	"sync"
	"time"
)

type item struct {
	value       *models.Order
	timeCreated time.Time
}

type shard struct {
	mutex sync.RWMutex
	items map[string]item
}

func newShard() shard {
	return shard{
		items: make(map[string]item),
	}
}

func (s *shard) set(key string, value *models.Order) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.items[key] = item{
		value:       value,
		timeCreated: time.Now(),
	}
}

func (s *shard) get(key string) (*models.Order, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	item, ex := s.items[key]
	if !ex {
		return nil, false
	}
	return item.value, ex
}

func (s *shard) countItems() uint64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return uint64(len(s.items))
}
