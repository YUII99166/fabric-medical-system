package cache

import (
	"encoding/json"
	"sync"
	"time"
)

// CacheItem 缓存项
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// MemoryCache 内存缓存
type MemoryCache struct {
	items map[string]*CacheItem
	mutex sync.RWMutex
}

// NewMemoryCache 创建新的内存缓存
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		items: make(map[string]*CacheItem),
	}
	
	// 启动定期清理过期缓存的goroutine
	go cache.cleanupExpired()
	
	return cache
}

// Set 设置缓存
func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items[key] = &CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(duration),
	}
}

// Get 获取缓存
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}
	
	// 检查是否过期
	if time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	
	return item.Data, true
}

// GetString 获取字符串缓存
func (c *MemoryCache) GetString(key string) (string, bool) {
	value, exists := c.Get(key)
	if !exists {
		return "", false
	}
	
	str, ok := value.(string)
	return str, ok
}

// GetJSON 获取JSON缓存并反序列化
func (c *MemoryCache) GetJSON(key string, target interface{}) bool {
	value, exists := c.Get(key)
	if !exists {
		return false
	}
	
	jsonStr, ok := value.(string)
	if !ok {
		return false
	}
	
	err := json.Unmarshal([]byte(jsonStr), target)
	return err == nil
}

// SetJSON 序列化为JSON并设置缓存
func (c *MemoryCache) SetJSON(key string, value interface{}, duration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	c.Set(key, string(jsonData), duration)
	return nil
}

// Delete 删除缓存
func (c *MemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.items, key)
}

// Clear 清空所有缓存
func (c *MemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items = make(map[string]*CacheItem)
}

// cleanupExpired 定期清理过期缓存
func (c *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.ExpiresAt) {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

// 全局缓存实例
var globalCache = NewMemoryCache()

// GetGlobalCache 获取全局缓存实例
func GetGlobalCache() *MemoryCache {
	return globalCache
}
