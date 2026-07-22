package main

import (
	"strings"
	"sync"
)

type DirCache struct {
	maxEntries  int
	mu          sync.RWMutex
	store       map[string]DirEntries
	previousDir string
}

func newDirCache(maxEntries int) *DirCache {
	return &DirCache{
		store:      make(map[string]DirEntries, maxEntries),
		maxEntries: maxEntries,
	}
}

func (c *DirCache) get(key string) (DirEntries, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entries, ok := c.store[key]; ok {
		return entries, true
	}

	return DirEntries{}, false
}

func (c *DirCache) set(key string, files DirEntries) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.store) >= c.maxEntries {
		tmp, ok := c.store[c.previousDir]
		c.store = make(map[string]DirEntries, c.maxEntries)

		if ok {
			c.store[c.previousDir] = tmp
		}
	}

	c.previousDir = key
	c.store[key] = files
}

func (c *DirCache) invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c *DirCache) invalidateRec(key string) {
	for k := range c.store {
		if strings.HasPrefix(k, key) || k == key {
			c.invalidate(k)
		}
	}
}
