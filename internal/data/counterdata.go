package data

import (
	"sync"
)

type myCounter struct {
	value counter
	mtx   sync.RWMutex
}

func NewCounter() *myCounter {
	return &myCounter{}
}

func (c *myCounter) IncCounter() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.value++
}

func (c *myCounter) Count() counter {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.value
}
