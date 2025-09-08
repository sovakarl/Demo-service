package memory

import (
	"time"
)

func (m *memory) startGC() {
	go m.gc()
}

func (m *memory) gc() {
	for index := range m.countShard {
		go func(i uint) {
			m.workGC(i)
		}(index)
	}
}

func (m *memory) workGC(indexShard uint) {
	ticker := time.NewTicker(m.defaultExpiration)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-m.chSignal:
			if !ok {
				return
			}
		case <-time.After(m.defaultExpiration):
			m.deleteAfterExpiration(indexShard)
		}
	}
}

func (c *memory) deleteAfterExpiration(indexShard uint) {
	shard := c.getShard(uint64(indexShard))
	sliceDeleteItem := []string{}

	shard.mutex.Lock()
	now := time.Now()
	for key, value := range shard.items {
		if now.Sub(value.timeCreated) >= c.defaultDuration {
			sliceDeleteItem = append(sliceDeleteItem, key)
		}
	}
	shard.mutex.Unlock()

	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	for _, key := range sliceDeleteItem {
		delete(shard.items, key)
	}
}

func (m *memory) stopGC() {
	close(m.chSignal)
}
