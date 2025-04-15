package main

import (
	"database/sql"
	"fmt"
	"http_server/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	env            string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	env := os.Getenv("ENV")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
	}
	databaseQueries := database.New(db)

	myApiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             databaseQueries,
		env:            env,
	}

	httpServerMux := http.NewServeMux()
	// App Endpoint
	fileServer := http.FileServer(http.Dir("."))
	httpServerMux.Handle("/app/", myApiConfig.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	// Health Endpoint
	httpServerMux.Handle("GET /api/healthz/", myApiConfig.middlewareMetricsInc(http.HandlerFunc(healthz)))
	httpServerMux.Handle("GET /admin/metrics", http.HandlerFunc(myApiConfig.metrics))
	httpServerMux.Handle("POST /admin/reset", http.HandlerFunc(myApiConfig.reset))
	httpServerMux.Handle("POST /api/validate_chirp", http.HandlerFunc(validate_chirp))
	httpServerMux.Handle("POST /api/users", http.HandlerFunc(myApiConfig.createUser))
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

func (cfg *apiConfig) metrics(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	hits := cfg.fileserverHits.Load()
	html, _ := os.ReadFile("metrics.html")
	text := fmt.Sprintf(string(html), hits)
	writer.Write([]byte(text))
}
