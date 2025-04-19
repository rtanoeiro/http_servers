package main

import (
	"database/sql"
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
	secret := os.Getenv("SECRET")
	api_key := os.Getenv("POLKA_APIKEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Println("Error connecting to the database:", err)
	}
	databaseQueries := database.New(db)

	myApiConfig := api.ApiConfig{
		FileserverHits: atomic.Int32{},
		Db:             databaseQueries,
		Env:            env,
		Secret:         secret,
		ApiKey:         api_key,
	}

	httpServerMux := http.NewServeMux()
	// App Endpoint
	fileServer := http.FileServer(http.Dir("."))
	httpServerMux.Handle("/app/", myApiConfig.MiddlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	// Health Endpoint
	httpServerMux.Handle("GET /api/healthz/", myApiConfig.MiddlewareMetricsInc(http.HandlerFunc(api.Healthz)))
	httpServerMux.Handle("GET /admin/metrics", http.HandlerFunc(myApiConfig.Metrics))
	httpServerMux.Handle("POST /admin/reset", http.HandlerFunc(myApiConfig.ResetUsers))
	httpServerMux.Handle("POST /api/login", http.HandlerFunc(myApiConfig.Login))

	httpServerMux.Handle("POST /api/users", http.HandlerFunc(myApiConfig.CreateUser))
	httpServerMux.Handle("PUT /api/users", http.HandlerFunc(myApiConfig.UpdateUser))
	httpServerMux.Handle("POST /api/polka/webhooks", http.HandlerFunc(myApiConfig.UpgradeUser))
	httpServerMux.Handle("POST /api/refresh", http.HandlerFunc(myApiConfig.Refresh))
	httpServerMux.Handle("POST /api/revoke", http.HandlerFunc(myApiConfig.Revoke))

	httpServerMux.Handle("DELETE /api/chirps/{chirpID}", http.HandlerFunc(myApiConfig.DeleteChirp))
	httpServerMux.Handle("POST /api/chirps", http.HandlerFunc(myApiConfig.InsertChirp))
	httpServerMux.Handle("GET /api/chirps", http.HandlerFunc(myApiConfig.GetChirps))
	httpServerMux.Handle("GET /api/chirps/{chirpID}", http.HandlerFunc(myApiConfig.GetSingleChirp))

	httpServer := http.Server{
		Handler: httpServerMux,
		Addr:    ":8080",
	}
	httpServer.ListenAndServe()
}
