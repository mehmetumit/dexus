package memcache

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mehmetumit/dexus/internal/core/ports"
)

type CacheMap map[string]string
type MemCache struct {
	sync.RWMutex //Useful for concurrent nonblocking reads
	logger       ports.Logger
	cacheMap     CacheMap
}

func NewMemCache(l ports.Logger) *MemCache {
	return &MemCache{
		logger:   l,
		cacheMap: make(CacheMap),
	}

}

func (mc *MemCache) expireAfter(key string, ttl time.Duration) {
	time.AfterFunc(ttl, func() {
		mc.Lock()
		defer mc.Unlock()
		delete(mc.cacheMap, key)
	})
}
func (mc *MemCache) Get(ctx context.Context, key string) (string, error) {
	mc.RLock()
	defer mc.RUnlock()
	val, ok := mc.cacheMap[key]
	if !ok {
		return "", ports.ErrKeyNotFound
	}
	return val, nil

}
func (mc *MemCache) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	mc.Lock()
	defer mc.Unlock()
	mc.cacheMap[key] = val
	mc.expireAfter(key, ttl)
	return nil
}

func (mc *MemCache) Delete(ctx context.Context, key string) error {
	mc.Lock()
	defer mc.Unlock()
	delete(mc.cacheMap, key)
	return nil
}
func (mc *MemCache) GenKey(ctx context.Context, s string) (string, error) {
	hashKey := uuid.NewSHA1(uuid.NameSpaceOID, []byte(s)).String()
	return hashKey, nil
}
func (mc *MemCache) Flush(ctx context.Context) error {
	mc.Lock()
	defer mc.Unlock()
	clear(mc.cacheMap)
	return nil
}
