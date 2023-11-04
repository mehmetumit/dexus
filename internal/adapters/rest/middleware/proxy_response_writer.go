package middleware

import (
	"net/http"
)

type ProxyResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	body       []byte
}

func NewProxyResponseWriter(w http.ResponseWriter) *ProxyResponseWriter {
	//Set default status code as 200
	return &ProxyResponseWriter{w, http.StatusOK, []byte{}}

}

func (pw *ProxyResponseWriter) Write(data []byte) (int, error) {
	pw.body = data
	return pw.ResponseWriter.Write(data)
}

// Proxy for getting status code of writer
func (pw *ProxyResponseWriter) WriteHeader(code int) {
	pw.ResponseWriter.WriteHeader(code)
	pw.StatusCode = code
}

//Header() Header -> Not implemented

func (pw *ProxyResponseWriter) GetLocation() string {
	return pw.Header().Get("location")
}
