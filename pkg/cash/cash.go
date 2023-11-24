package cash

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	elements          map[string]Element
}

type Element struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func NewCash(defaultExpiration, cleanupInterval time.Duration) *Cache {
	elem := make(map[string]Element)
	cash := &Cache{
		elements:          elem,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}
	if cleanupInterval > 0 {
		go cash.StartGC()
	}
	return cash
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	element, ok := c.elements[key]
	if !ok {
		return nil, false
	}
	if element.Expiration > 0 {
		if time.Now().UnixNano() > element.Expiration {
			return nil, false
		}
	}
	return element.Value, true
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64
	if duration == 0 {
		duration = c.defaultExpiration
	}
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	c.Lock()
	defer c.Unlock()
	c.elements[key] = Element{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()
	if _, found := c.elements[key]; !found {
		return errors.New("Key not found")
	}
	delete(c.elements, key)
	return nil
}

func (c *Cache) StartGC() {
	go c.GC()
}

func (c *Cache) GC() {
	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)
		if c.elements == nil {
			return
		}
		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearElements(keys)
		}
	}

}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()
	for k, i := range c.elements {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}
	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearElements(keys []string) {
	c.Lock()
	defer c.Unlock()
	for _, k := range keys {
		delete(c.elements, k)
	}
}
