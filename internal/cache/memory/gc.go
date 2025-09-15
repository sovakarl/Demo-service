package memory

import (
	"errors"
	"log/slog"
	"sync"
	"time"
)

var (
	ErrGCStopped    = errors.New("gc stoped")
	ErrGCWasStopped = errors.New("GC was stopped")
)



func (m *memory) startGC() {
	go m.gc()
}

func (m *memory) gc() {
	m.logger.Debug("start cache GC")

	var wg sync.WaitGroup
	wg.Add(int(m.countShard))

	for indexShard := range m.countShard {
		logger := m.logger.With("shard", indexShard)

		go func(i uint64, log *slog.Logger) {
			defer wg.Done()
			m.workGC(i, log)
		}(indexShard, logger)

	}
	wg.Wait()
	m.chStatusGC <- struct{}{}
}

func (m *memory) workGC(indexShard uint64, logger *slog.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Debug("GC panic", "panic", r)
		}
	}()

	shard := m.getShard(indexShard)
	ticker := time.NewTicker(m.defaultExpiration)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-m.chSignal:
			if !ok {
				logger.Debug("shard stoped")
				return
			}
		case <-ticker.C:
			if shard.countItems() == 0 {
				continue
			}
			logger.Debug("GC start clean ")
			m.deleteAfterExpiration(shard, logger)
		}
	}
}

func (m *memory) deleteAfterExpiration(shard *shard, logger *slog.Logger) {

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
	logger.Debug("GC clean", "count delete", rune(len(sliceDeleteItem)))
}

func (m *memory) stopGC() error {
	m.logger.Debug("get signal for stop GC")

	select {
	case _, ex := <-m.chSignal:
		if !ex {
			return ErrGCWasStopped
		}
		close(m.chSignal)
	default:
		close(m.chSignal)
	}

	<-m.chStatusGC
	return nil

}
