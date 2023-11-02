package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mehmetumit/dexus/internal/core/ports"
)

type MockCacheMap map[string]string
type MockCacher struct {
	sync.RWMutex //Useful for concurrent nonblocking reads
	cacheMap     MockCacheMap
}

func NewMockCacher() *MockCacher {
	return &MockCacher{
		cacheMap: make(MockCacheMap),
	}
}

func (mc *MockCacher) expireAfter(key string, ttl time.Duration) {
	time.AfterFunc(ttl, func() {
		mc.Lock()
		defer mc.Unlock()
		delete(mc.cacheMap, key)
	})
}
func (mc *MockCacher) GenKey(ctx context.Context, s string) (string, error) {
	hashKey := uuid.NewSHA1(uuid.NameSpaceOID, []byte(s)).String()
	return hashKey, nil

}
func (mc *MockCacher) Get(ctx context.Context, key string) (string, error) {
	mc.RLock()
	defer mc.RUnlock()
	val, ok := mc.cacheMap[key]
	if !ok {
		return "", ports.ErrKeyNotFound
	}
	return val, nil

}
func (mc *MockCacher) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	mc.Lock()
	defer mc.Unlock()
	mc.cacheMap[key] = val
	mc.expireAfter(key, ttl)
	return nil
}
func (mc *MockCacher) Delete(ctx context.Context, key string) error {
	mc.Lock()
	defer mc.Unlock()
	delete(mc.cacheMap, key)
	return nil
}
func (mc *MockCacher) Flush(ctx context.Context) error {
	mc.Lock()
	defer mc.Unlock()
	clear(mc.cacheMap)
	return nil
}
