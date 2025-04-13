package main

import (
	"net/http"
)

func main() {
	httpServerMux := http.NewServeMux()
	// App Endpoint
	fileServer := http.FileServer(http.Dir("."))
	httpServerMux.Handle("/app/", http.StripPrefix("/app", fileServer))
	// Health Endpoint
	httpServerMux.HandleFunc("/healthz/", healthz)

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
