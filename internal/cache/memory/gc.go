package memory

import (
	"time"
)

func (m *memory) startGC() {
	go m.gc()
}

func (m *memory) gc() {
	for index := range m.countShard {
		go func(i uint64) {
			m.workGC(i)
		}(index)
	}
}

func (m *memory) workGC(indexShard uint64) {
	defer func() {
		if r := recover(); r != nil {
			m.logger.Debug("GC panic", "shard", indexShard, "panic", r)
		}
	}()
	shard := m.getShard(indexShard)
	ticker := time.NewTicker(m.defaultExpiration)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-m.chSignal:
			if !ok {
				m.logger.Debug("GC stoped ", "shard", indexShard)
				return
			}
		case <-ticker.C:
			if shard.countItems() == 0 {
				continue
			}
			m.logger.Debug("GC start clean ", "shard", indexShard)
			m.deleteAfterExpiration(shard)
		}
	}
}

func (m *memory) deleteAfterExpiration(shard *shard) {

	sliceDeleteItem := []string{}

	shard.mutex.RLock()
	now := time.Now()
	for key, value := range shard.items {
		if now.Sub(value.timeCreated) >= m.defaultDuration {
			sliceDeleteItem = append(sliceDeleteItem, key)
		}
	}
	shard.mutex.RUnlock()

	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	for _, key := range sliceDeleteItem {
		delete(shard.items, key)
	}
	m.logger.Debug("GC clean", "indexSh", "count delete", len(sliceDeleteItem))
}

func (m *memory) stopGC() error {
	m.logger.Debug("get signal for stop GC")
	if _, ex := <-m.chSignal; ex {
		close(m.chSignal)
		return nil
	}
	return nil
}
