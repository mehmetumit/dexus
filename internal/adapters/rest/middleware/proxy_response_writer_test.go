package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestResponseWriter(t testing.TB, w http.ResponseWriter) *ProxyResponseWriter {
	t.Helper()
	return NewProxyResponseWriter(w)

}
func TestProxyResponseWriter_Write(t *testing.T) {
	testDatas := []string{
		"test1",
		"",
		"<a href=\"https://foo.com\">Found</a>",
		"  ",
	}
	for _, data := range testDatas {
		recorder := httptest.NewRecorder()
		proxyWriter := newTestResponseWriter(t, recorder)
		proxyWriter.Write([]byte(data))
		if !bytes.Equal(recorder.Body.Bytes(), []byte(data)) {
			t.Errorf("expected body data \"%s\", got \"%s\"", data, recorder.Body.Bytes())
		}
	}

}
func TestProxyResponseWriter_WriteHeader(t *testing.T) {
	statusCodes := []int{
		http.StatusAccepted,
		http.StatusBadGateway,
		http.StatusNotFound,
		http.StatusFound,
		http.StatusPermanentRedirect,
	}
	for _, code := range statusCodes {

		recorder := httptest.NewRecorder()
		proxyWriter := newTestResponseWriter(t, recorder)
		proxyWriter.WriteHeader(code)
		if code != recorder.Code {
			t.Errorf("exprected status code %d, got %d", code, recorder.Code)
		}
	}

}
func TestProxyResponseWriter_GetLocation(t *testing.T) {
	locations := []string{
		"https://foo.com",
		"http://bar.com",
		"https://foo-bar.com",
		"https://foo.bar.com",
		"https://foo.bar.com",
		"https://foo.bar.com/foo/bar",
		"https://foo.bar.com/foo/bar?foo=bar&bar=foo",
	}
	for _, loc := range locations {
		recorder := httptest.NewRecorder()
		mockReq := httptest.NewRequest("GET", "https://redirect.example.com/go", nil)
		proxyWriter := newTestResponseWriter(t, recorder)

		http.Redirect(proxyWriter, mockReq, loc, http.StatusFound)
		gotLocation := proxyWriter.GetLocation()
		if loc != gotLocation {
			t.Errorf("expected redirection location %v, got %v", loc, gotLocation)
		}
	}

}
