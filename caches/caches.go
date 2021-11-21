package caches

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	data    map[string]*value // byte is net trans ok
	options Options
	status  Status
	lock    *sync.RWMutex //read and write lock
}

func NewCache() *Cache {
	return NewCacheWith(*DefaultOptions())
}

func NewCacheWith(options Options) *Cache {
	return &Cache{
		data:    make(map[string]*value, 256),
		options: options,
		status:  *NewStauts(),
		lock:    &sync.RWMutex{},
	}
}
func (c *Cache) Get(key string) ([]byte, bool) {
	//read lock
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	if !ok {
		return nil, false
	}
	if !value.alive() {
		c.lock.RUnlock()
		c.Delete(key)
		c.lock.RLock()
		return nil, false
	}
	return value.visit(), true // lru ttl expire mechanism
}

func (c *Cache) Set(key string, value []byte) error {
	return c.SetWithTTL(key, value, NeverDie)
}

func (c *Cache) SetWithTTL(key string, value []byte, ttl int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.data)
	}
	if !c.checkEntrySize(key, value) {
		if oldValue, ok := c.data[key]; ok {
			c.status.addEntry(key, oldValue.data)
		}
		return errors.New("the entry size will exceed if you set this entry")
	}
	c.status.addEntry(key, value)
	c.data[key] = newValue(value, ttl)
	return nil
}
func (c *Cache) Delete(key string) {
	//write lock
	c.lock.Lock()
	defer c.lock.Unlock()
	if oldValue, ok := c.data[key]; ok {
		c.status.subEntry(key, oldValue.data)
		delete(c.data, key)
	}
}

func (c *Cache) Status() Status {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return *&c.status
}

func (c *Cache) checkEntrySize(newKey string, newValue []byte) bool {
	return c.status.entrySize()+int64(len(newKey))+int64(len(newValue)) <= c.options.MaxEntrySize*1024*1024
}

//clean expire data
func (c *Cache) gc() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// use count record current clean nums
	count := 0
	for key, value := range c.data {
		if !value.alive() {
			c.status.subEntry(key, value.data)
			delete(c.data, key)
			count++
			if count >= int(c.options.MaxGCCount) {
				break
			}
		}
	}
}

func (c *Cache) AutoGc() {
	go func() {
		ticker := time.NewTicker(time.Duration(c.options.GCDuration) * time.Minute)
		for range ticker.C {
			c.gc()
		}
	}()
}
