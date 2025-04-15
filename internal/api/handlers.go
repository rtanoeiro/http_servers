package api

import (
	"context"
	"encoding/json"
	"fmt"
	"http_server/internal/database"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func (config *ApiConfig) CreateUser(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	user := UserAdd{}
	err := decoder.Decode(&user)

	if err != nil {
		errorMsg := ChirpMsgError{
			Error: "Something went wrong",
		}
		data, err := json.Marshal(errorMsg)

		if err != nil {
			respondWithError(writer, http.StatusInternalServerError, "Error marshalling JSON during initial request check")
			return
		}
		respondWithJSON(writer, http.StatusOK, data)
		return
	}

	createUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     user.Email,
	}
	createdUser, createError := config.Db.CreateUser(request.Context(), createUser)
	if createError != nil {
		log.Println("CreateUser error:", createError)
		respondWithError(writer, http.StatusInternalServerError, "Unable to create user")
		return
	}

	returnUser := database.User{
		ID:        createUser.ID,
		CreatedAt: createUser.CreatedAt,
		UpdatedAt: createUser.UpdatedAt,
		Email:     createdUser.Email,
	}
	data, err := json.Marshal(returnUser)

	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error marshalling JSON during initial request check")
		return
	}
	respondWithJSON(writer, http.StatusCreated, data)

}

func (cfg *ApiConfig) Reset(writer http.ResponseWriter, request *http.Request) {
	if cfg.Env == "dev" {
		deleteError := cfg.Db.DeleteAllUsers(context.Background())

		if deleteError != nil {
			respondWithError(writer, http.StatusInternalServerError, "Unable to delete all users")
		}
	} else {
		respondWithError(writer, http.StatusForbidden, "Unable to perform this action in this environment")
	}
}

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
