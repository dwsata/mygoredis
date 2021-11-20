package caches

import (
	"github.com/opentechnologysel/mygoredis/helpers"
	"sync"
)

type Cache struct {
	data  map[string][]byte // byte is net trans ok
	lock  *sync.RWMutex     //read and write lock
	count int64             // count key nums
}

func NewCache() *Cache {
	return &Cache{
		data:  make(map[string][]byte),
		lock:  &sync.RWMutex{},
		count: 0,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	//read lock
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Set(key string, value []byte) {
	//write lock
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	if !ok {
		c.count++
	}
	c.data[key] = helpers.Copy(value)

}

func (c *Cache) Delete(key string) {
	//write lock
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	if ok {
		c.count--
		delete(c.data, key)
	}
}

func (c *Cache) Count() int64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.count
}
