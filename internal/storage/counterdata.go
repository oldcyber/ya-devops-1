package storage

import (
	"sync"
)

type MyCounter struct {
	name  string
	value Counter
	mtx   sync.RWMutex
}

func NewCounter() *MyCounter {
	return &MyCounter{}
}

func (c *MyCounter) IncCounter() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.name = "PollCount"
	c.value++
}

func (c *MyCounter) Count() Counter {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.value
}
