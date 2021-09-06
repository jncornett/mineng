// Package cache provides a resettable variant of sync.Once.
package cache

import (
	"sync"
)

// Cache is a resettable variant of sync.Once.
type Cache struct {
	o *sync.Once
	// m provides mutual exclusion for readers of the sync.Once pointer o,
	// and writers of the same pointer.
	m sync.RWMutex
	// i provides one-time initialization of the sync.Once object in o.
	i sync.Once
}

func (c *Cache) Do(fn func()) {
	c.m.RLock() // exclude calls to c.Reset() for the duration of this function.
	defer c.m.RUnlock()
	// we excluded callers of c.Reset(), but not other callers of c.Do()
	// gate 1st time initialization of c.o here.
	c.i.Do(func() { c.o = new(sync.Once) })
	c.o.Do(fn)
}

func (c *Cache) Reset() {
	c.m.Lock() // exclude calls to c.Do() for the duration of this function.
	defer c.m.Unlock()
	// no need to gate with c.i, because at most one caller can acquire the write lock on the mutex.
	c.o = new(sync.Once)
}
