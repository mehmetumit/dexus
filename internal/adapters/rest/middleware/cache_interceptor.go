package middleware

import (
	"net/http"
	"time"

	"github.com/mehmetumit/dexus/internal/core/ports"
)


type CacheInterceptor struct {
	cacher ports.Cacher
	logger ports.Logger
	ttl time.Duration
}

func NewCacheInterceptor(c ports.Cacher, l ports.Logger, ttl time.Duration) CacheInterceptor {
	return CacheInterceptor{
		cacher: c,
		logger: l,
		ttl: ttl,
	}

}
func (ch *CacheInterceptor) InterceptHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxyWriter := NewProxyResponseWriter(w)
		defer func() {
			// Don't cache observability paths and some http methods
			if r.URL.Path == "/health" ||
				r.URL.Path == "/metrics" ||
				r.URL.Path == "/monitor" ||
				r.Method == http.MethodOptions || r.Method == http.MethodPost {
				ch.logger.Debug("Pass caching")
				next.ServeHTTP(w, r)
				return
			}
			ctx := r.Context()

			hashURL, _ := ch.cacher.GenKey(ctx, r.URL.Path)
			// The method is not a query, which means a state change occurs on result of this URL
			if r.Method != http.MethodGet {
				//Invalidate cache
				next.ServeHTTP(w, r)
				err := ch.cacher.Delete(ctx, hashURL)
				ch.logger.Debug("Invalidate cache:", hashURL)
				if err != nil {
					ch.logger.Error("command http method cache invalidation err:", err)
				}
				return
			}
			cacheData, err := ch.cacher.Get(ctx, hashURL)

			if err != nil || len(cacheData) == 0 {
				ch.logger.Debug("Cache miss:", hashURL)
				if err != ports.ErrKeyNotFound{
					ch.logger.Error("internal cache error:",err)
				}
				// Get response of request by sending it to next
				next.ServeHTTP(proxyWriter, r)
				if proxyWriter.StatusCode == http.StatusFound {
					// Set cache using redirection location which stored in proxyWriter
					ch.cacher.Set(ctx, hashURL, proxyWriter.GetLocation(), ch.ttl)
				}
				return// Response already built using proxyWriter
			}
			// Don't send the request to the next
			//Instead, respond to the client with cached data
			ch.logger.Debug("Cache hit:", hashURL)
			ch.logger.Debug("Cached data:", cacheData)
			w.Header().Add("x-cached-response", "true")

			http.Redirect(w, r, cacheData, http.StatusFound)
			return
		}()
	})
}
