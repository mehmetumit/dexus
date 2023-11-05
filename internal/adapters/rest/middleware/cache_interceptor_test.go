package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mehmetumit/dexus/internal/mocks"
)

func newTestCacheInterceptor(tb testing.TB, ttl time.Duration) CacheInterceptor {
	tb.Helper()
	c := mocks.NewMockCacher()
	l := mocks.NewMockLogger()
	return NewCacheInterceptor(c, l, ttl)
}

type passInterceptorConfig struct {
	interceptor    *CacheInterceptor
	reqPath        string
	reqMethod      string
	expectedURL    string
	expectedStatus int
}

func genSinglePassInterceptorResult(t *testing.T, cfg passInterceptorConfig) *httptest.ResponseRecorder {
	t.Helper()
	mockResp := httptest.NewRecorder()
	mockReq := httptest.NewRequest(cfg.reqMethod, "https://foo.com"+cfg.reqPath, nil)
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, string(cfg.expectedURL), cfg.expectedStatus)
	})

	interceptHandler := cfg.interceptor.InterceptHandler(nextHandler)
	interceptHandler.ServeHTTP(mockResp, mockReq)
	return mockResp
}
func genDoublePassInterceptorResult(t *testing.T, cfg passInterceptorConfig) *httptest.ResponseRecorder {
	t.Helper()
	// Run twice to check caching
	// Return only second response
	genSinglePassInterceptorResult(t, cfg)
	return genSinglePassInterceptorResult(t, cfg)
}

type chechCachedConfig struct {
	reqPath        string
	expectedURL    string
	expectedStatus int
}

func errorIfCached(t *testing.T, mockResp *httptest.ResponseRecorder, cfg chechCachedConfig) {
	t.Helper()
	if mockResp.Header().Get("x-cached-response") != "" {
		t.Errorf("path %s must not be cached", cfg.reqPath)
	}
	if mockResp.Header().Get("location") != cfg.expectedURL {
		t.Errorf("expected location data %s, got %s", cfg.expectedURL, mockResp.Header().Get("location"))
	}
	if mockResp.Code != cfg.expectedStatus {
		t.Errorf("expected status code %d, got %d", cfg.expectedStatus, mockResp.Code)
	}

}
func errorIfNotCached(t *testing.T, mockResp *httptest.ResponseRecorder, cfg chechCachedConfig) {
	t.Helper()
	if mockResp.Header().Get("x-cached-response") != "true" {
		t.Errorf("path %s must be cached", cfg.reqPath)
	}
	if mockResp.Header().Get("location") != cfg.expectedURL {
		t.Errorf("expected location data %s, got %s", cfg.expectedURL, mockResp.Header().Get("location"))
	}
	if mockResp.Code != cfg.expectedStatus {
		t.Errorf("expected status code %d, got %d", cfg.expectedStatus, mockResp.Code)
	}
}
func TestCacheInterceptor_DontCachePaths(t *testing.T) {
	interceptor := newTestCacheInterceptor(t, 20*time.Second)
	paths := []string{
		"/metrics",
		"/health",
		"/monitor",
	}
	testURL := "https://foor-bar.example.com/abc?foo=bar&bar=foo"
	testStatus := http.StatusAccepted
	for _, path := range paths {

		mockResp := genDoublePassInterceptorResult(t, passInterceptorConfig{
			interceptor:    &interceptor,
			reqPath:        path,
			reqMethod:      "GET",
			expectedURL:    testURL,
			expectedStatus: testStatus,
		})

		errorIfCached(t, mockResp, chechCachedConfig{
			reqPath:        path,
			expectedURL:    testURL,
			expectedStatus: testStatus,
		})
	}
}

func TestCacheInterceptor_DontCacheReqMethods(t *testing.T) {
	interceptor := newTestCacheInterceptor(t, 20*time.Second)
	methods := []string{
		"OPTIONS",
		"POST",
	}
	path := "/test"
	testURL := "https://foor-bar.example.com/abc?foo=bar&bar=foo"
	testStatus := http.StatusOK

	for _, m := range methods {

		mockResp := genDoublePassInterceptorResult(t, passInterceptorConfig{
			interceptor:    &interceptor,
			reqPath:        path,
			reqMethod:      m,
			expectedURL:    testURL,
			expectedStatus: testStatus,
		})

		errorIfCached(t, mockResp, chechCachedConfig{
			reqPath:        path,
			expectedURL:    testURL,
			expectedStatus: testStatus,
		})
	}

}

func TestCacheInterceptor_CacheInvalidation(t *testing.T) {
	testURL := "https://foor-bar.example.com/abc?foo=bar&bar=foo"
	newURL := "https://new-foo-bar.example.com/new?new=true"
	expectedStatus := http.StatusFound
	path := "/test"
	t.Run("DELETE method cache invalidation", func(t *testing.T) {
		interceptor := newTestCacheInterceptor(t, 20*time.Second)

		cachedResp := genDoublePassInterceptorResult(t, passInterceptorConfig{
			interceptor:    &interceptor,
			reqPath:        path,
			reqMethod:      "GET",
			expectedURL:    testURL,
			expectedStatus: expectedStatus,
		})
		errorIfNotCached(t, cachedResp, chechCachedConfig{
			reqPath:        path,
			expectedURL:    testURL,
			expectedStatus: expectedStatus,
		})
		invalidatedResp := genSinglePassInterceptorResult(t, passInterceptorConfig{
			interceptor:    &interceptor,
			reqPath:        path,
			reqMethod:      "DELETE",
			expectedURL:    newURL,
			expectedStatus: expectedStatus,
		})
		errorIfCached(t, invalidatedResp, chechCachedConfig{
			reqPath:        path,
			expectedURL:    newURL,
			expectedStatus: expectedStatus,
		})
		freshResp := genSinglePassInterceptorResult(t, passInterceptorConfig{
			interceptor:    &interceptor,
			reqPath:        path,
			reqMethod:      "GET",
			expectedURL:    newURL,
			expectedStatus: expectedStatus,
		})
		errorIfCached(t, freshResp, chechCachedConfig{
			reqPath:        path,
			expectedURL:    newURL,
			expectedStatus: expectedStatus,
		})
	})
	t.Run("PUT method cache invalidation", func(t *testing.T) {
		interceptor := newTestCacheInterceptor(t, 20*time.Second)
		cachedResp := genDoublePassInterceptorResult(t, passInterceptorConfig{
			interceptor:    &interceptor,
			reqPath:        path,
			reqMethod:      "GET",
			expectedURL:    testURL,
			expectedStatus: expectedStatus,
		})
		errorIfNotCached(t, cachedResp, chechCachedConfig{
			reqPath:        path,
			expectedURL:    testURL,
			expectedStatus: expectedStatus,
		})
		invalidatedResp := genSinglePassInterceptorResult(t, passInterceptorConfig{
			interceptor:    &interceptor,
			reqPath:        path,
			reqMethod:      "PUT",
			expectedURL:    newURL,
			expectedStatus: expectedStatus,
		})
		errorIfCached(t, invalidatedResp, chechCachedConfig{
			reqPath:        path,
			expectedURL:    newURL,
			expectedStatus: expectedStatus,
		})
		freshResp := genSinglePassInterceptorResult(t, passInterceptorConfig{
			interceptor:    &interceptor,
			reqPath:        path,
			reqMethod:      "GET",
			expectedURL:    newURL,
			expectedStatus: expectedStatus,
		})
		errorIfCached(t, freshResp, chechCachedConfig{
			reqPath:        path,
			expectedURL:    newURL,
			expectedStatus: expectedStatus,
		})
	})

}
func TestCacheInterceptor_DontCacheNotFound(t *testing.T) {
	interceptor := newTestCacheInterceptor(t, 20*time.Second)
	testURL := "https://foor-bar.example.com/abc?foo=bar&bar=foo"
	expectedStatus := http.StatusNotFound
	path := "/test"

	notFoundResp := genDoublePassInterceptorResult(t, passInterceptorConfig{
		interceptor:    &interceptor,
		reqPath:        path,
		reqMethod:      "GET",
		expectedURL:    testURL,
		expectedStatus: expectedStatus,
	})
	errorIfCached(t, notFoundResp, chechCachedConfig{
		reqPath:        path,
		expectedURL:    testURL,
		expectedStatus: expectedStatus,
	})
}
func TestCacheInterceptor_TTLInvalidation(t *testing.T) {
	ttl := 50 * time.Millisecond
	interceptor := newTestCacheInterceptor(t, ttl)
	testURL := "https://foor-bar.example.com/abc?foo=bar&bar=foo"
	expectedStatus := http.StatusFound
	path := "/test"

	cachedResp := genDoublePassInterceptorResult(t, passInterceptorConfig{
		interceptor:    &interceptor,
		reqPath:        path,
		reqMethod:      "GET",
		expectedURL:    testURL,
		expectedStatus: expectedStatus,
	})
	errorIfNotCached(t, cachedResp, chechCachedConfig{
		reqPath:        path,
		expectedURL:    testURL,
		expectedStatus: expectedStatus,
	})
	time.Sleep(ttl * 2)
	expiredCacheResp := genSinglePassInterceptorResult(t, passInterceptorConfig{
		interceptor:    &interceptor,
		reqPath:        path,
		reqMethod:      "GET",
		expectedURL:    testURL,
		expectedStatus: expectedStatus,
	})

	errorIfCached(t, expiredCacheResp, chechCachedConfig{
		reqPath:        path,
		expectedURL:    testURL,
		expectedStatus: expectedStatus,
	})
}
