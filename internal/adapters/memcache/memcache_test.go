package memcache

import (
	"context"
	"testing"
	"time"

	"github.com/mehmetumit/dexus/internal/core/ports"
	"github.com/mehmetumit/dexus/internal/mocks"
)

func newTestMemCache(t testing.TB) *MemCache {
	t.Helper()
	mockLogger := mocks.NewMockLogger()
	return NewMemCache(mockLogger)
}

func TestMemCache_GenKey_Set_Get_GetNotFound_Flush_ExpireAfter(t *testing.T) {
	memCache := newTestMemCache(t)
	ctx := context.Background()
	cacheTTL := 5 * time.Second
	cacheMap := map[string]string{
		"a-test-path":  "https://test-path.com",
		"this/is/test": "https://this-is-test.com",
		"test1234":     "http://test1234.com",
	}
	//Hashed key : path key
	var keyMap map[string]string = make(map[string]string)

	t.Run("GenKey Set Get", func(t *testing.T) {
		for k, v := range cacheMap {
			hashKey, err := memCache.GenKey(ctx, k)
			if err != nil {
				t.Errorf("Expected error nil, got %v", err)
			}
			keyMap[hashKey] = k
			err = memCache.Set(ctx, hashKey, v, cacheTTL)
			if err != nil {
				t.Errorf("Expected error nil, got %v", err)
			}

		}
		for k, v := range keyMap {
			val, err := memCache.Get(ctx, k)
			if err != nil {
				t.Errorf("Expected error nil, got %v", err)
			}
			if val != cacheMap[v] {
				t.Errorf("Expected redirection url %v, got %v", cacheMap[v], val)
			}
		}
	})

	t.Run("Get Not Found", func(t *testing.T) {

		notFoundKeys := []string{
			"not-found",
			"1234/not/found",
			"",
			" ",
		}
		for _, nfk := range notFoundKeys {
			_, err := memCache.Get(ctx, nfk)
			if err != ports.ErrKeyNotFound {
				t.Errorf("Expected error %v, got %v", ports.ErrKeyNotFound, err)
			}
		}
	})
	t.Run("Flush", func(t *testing.T) {

		err := memCache.Flush(ctx)
		if err != nil {
			t.Errorf("Expected error nil, got %v", err)
		}
		for k := range keyMap {
			_, err := memCache.Get(ctx, k)
			if err != ports.ErrKeyNotFound {
				t.Errorf("Expected error %v, got %v", ports.ErrKeyNotFound, err)
			}

		}
	})
	t.Run("Delete", func(t *testing.T) {
		for k, v := range cacheMap {
			hashKey, err := memCache.GenKey(ctx, k)
			if err != nil {
				t.Errorf("Expected error nil, got %v", err)
			}
			err = memCache.Set(ctx, hashKey, v, cacheTTL)
			if err != nil {
				t.Errorf("Expected error nil, got %v", err)
			}
		}
		for k := range keyMap {
			err := memCache.Delete(ctx, k)
			if err != nil {
				t.Errorf("Expected error nil, got %v", err)
			}
			_, err = memCache.Get(ctx, k)
			if err != ports.ErrKeyNotFound {
				t.Errorf("Expected error %v, got %v", ports.ErrKeyNotFound, err)
			}
		}

	})
	t.Run("Expire After", func(t *testing.T) {
		ttl := 50 * time.Millisecond
		for k, v := range cacheMap {
			hashKey, err := memCache.GenKey(ctx, k)
			if err != nil {
				t.Errorf("Expected error nil, got %v", err)
			}
			err = memCache.Set(ctx, hashKey, v, ttl)
			if err != nil {
				t.Errorf("Expected error nil, got %v", err)
			}
		}

		time.Sleep(ttl * 2)
		for k := range keyMap {
			_, err := memCache.Get(ctx, k)
			if err != ports.ErrKeyNotFound {
				t.Errorf("Expected error %v, got %v", ports.ErrKeyNotFound, err)
			}
		}
	})

}
