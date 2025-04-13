package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	myApiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	httpServerMux := http.NewServeMux()
	// App Endpoint
	fileServer := http.FileServer(http.Dir("."))
	httpServerMux.Handle("/app/", myApiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	// Health Endpoint
	httpServerMux.Handle("/healthz/", myApiConfig.middlewareMetricsInc(http.HandlerFunc(healthz)))
	httpServerMux.Handle("/metrics", http.HandlerFunc(myApiConfig.metrics))
	httpServerMux.Handle("/reset", http.HandlerFunc(myApiConfig.reset))

	httpServer := http.Server{
		Handler: httpServerMux,
		Addr:    ":8080",
	}
	httpServer.ListenAndServe()
}

func healthz(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cfg.fileserverHits.Add(1)
		// This automatically starts the next handlers after incrementing hit count
		next.ServeHTTP(writer, request)
	})
}

func (cfg *apiConfig) reset(writer http.ResponseWriter, request *http.Request) {
	cfg.fileserverHits = atomic.Int32{}
	header := writer.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Server hits reduced to 0"))
}

func (cfg *apiConfig) metrics(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	header.Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	text := "Hits: " + fmt.Sprintf("%d", hits)
	writer.Write([]byte(text))
}
