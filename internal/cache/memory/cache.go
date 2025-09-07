package memory

import (
	"demo-service/internal/models"
	"sync"
	"time"
)

// type logger struct{}

// func (l logger)newSet(item Item)string{
// 	fmt.Sprint("New set %v ,all count %v ",item)
// }

type Item struct {
	value       *models.Order
	timeCreated time.Time
}

// type infoCache struct {
// }

type memory struct {
	flagLog           bool
	countItem         int64
	mutex             sync.Mutex
	chSignal          chan struct{}
	chLog             chan string
	items             map[string]Item
	defaultExpiration time.Duration
	defaultDuration   time.Duration
}

func NewCache(defaultExpiration, defaultDuration time.Duration) *memory {
	items := make(map[string]Item)
	cache := memory{
		flagLog:           false,
		chLog:             nil,
		chSignal:          make(chan struct{}),
		defaultExpiration: defaultExpiration,
		defaultDuration:   defaultDuration,
		items:             items,
	}
	if defaultDuration > 0 {
		cache.startGC()
	}
	return &cache
}
func (c *memory) Set(order *models.Order) {
	orderUID := order.GetUid()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ex := c.items[orderUID]; !ex {
		c.countItem++
	}

	c.items[orderUID] = Item{
		timeCreated: time.Now(),
		value:       order,
	}
}

func (c *memory) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.del(key)
}

func (c *memory) Get(key string) (*models.Order, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	item, ex := c.items[key]
	if !ex {
		return nil, false
	}
	return item.value, true
}

func (c *memory) del(key string) {
	if _, ex := c.items[key]; ex {
		c.countItem--
		delete(c.items, key)
	}
}

func (c *memory) Logs() chan string {
	if !c.flagLog {
		c.flagLog = true
		c.chLog = make(chan string)
		return c.chLog
	}
	return nil
}

func (c *memory) startGC() {
	go c.gc()
}

func (c *memory) Close() {
	c.stopGC()
	c.stopLog()
}

func (c *memory) gc() {
	for {
		select {
		case <-c.chSignal:
			return
		case <-time.After(c.defaultExpiration):
			c.deleteAfterExpiration()
		}
	}
}

func (c *memory) deleteAfterExpiration() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, value := range c.items {
		now := time.Now()

		if now.Sub(value.timeCreated) >= c.defaultDuration {
			c.del(key)
		}
	}
}

func (m *memory) stopGC() {
	m.chSignal <- struct{}{}
}

func (m *memory) stopLog() {
	if m.flagLog {
		close(m.chLog)
	}
}
