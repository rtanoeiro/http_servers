package api

import (
	"fmt"

	"net/http"
	"os"
)

// Get Methods
func Healthz(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.FileserverHits.Add(1)
		// This automatically starts the next handlers after incrementing hit count
		next.ServeHTTP(writer, request)
	})
}

func (cfg *ApiConfig) Metrics(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	hits := cfg.FileserverHits.Load()
	html, _ := os.ReadFile("metrics.html")
	text := fmt.Sprintf(string(html), hits)
	writer.Write([]byte(text))
}
