package main

import (
	"database/sql"
	"fmt"
	"http_server/internal/api"
	"http_server/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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

	myApiConfig := api.ApiConfig{
		FileserverHits: atomic.Int32{},
		Db:             databaseQueries,
		Env:            env,
	}

	httpServerMux := http.NewServeMux()
	// App Endpoint
	fileServer := http.FileServer(http.Dir("."))
	httpServerMux.Handle("/app/", myApiConfig.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	// Health Endpoint
	httpServerMux.Handle("GET /api/healthz/", myApiConfig.MiddlewareMetricsInc(http.HandlerFunc(api.Healthz)))
	httpServerMux.Handle("GET /admin/metrics", http.HandlerFunc(myApiConfig.Metrics))
	httpServerMux.Handle("POST /admin/reset", http.HandlerFunc(myApiConfig.Reset))
	httpServerMux.Handle("POST /api/users", http.HandlerFunc(myApiConfig.CreateUser))
	httpServerMux.Handle("POST /api/chirps", http.HandlerFunc(myApiConfig.InsertChirp))
	httpServer := http.Server{
		Handler: httpServerMux,
		Addr:    ":8080",
	}
	httpServer.ListenAndServe()
}
