package main

import (
	"strings"
	"sync"

	goadb "github.com/electricbubble/gadb"
)

type DirCache struct {
	maxEntries  int
	mu          sync.RWMutex
	store       map[string][]goadb.DeviceFileInfo
	previousDir string
}

func newDirCache(maxEntries int) *DirCache {
	return &DirCache{
		store:      make(map[string][]goadb.DeviceFileInfo, maxEntries),
		maxEntries: maxEntries,
	}
}

func (c *DirCache) get(key string) ([]goadb.DeviceFileInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if files, ok := c.store[key]; ok {
		return files, true
	}

	return nil, false
}

func (c *DirCache) set(key string, files []goadb.DeviceFileInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.store) >= c.maxEntries {
		tmp, ok := c.store[c.previousDir]
		c.store = make(map[string][]goadb.DeviceFileInfo, c.maxEntries)

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
