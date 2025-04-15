package main

import (
	"context"
	"encoding/json"
	"http_server/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type chirpMsg struct {
	Body string `json:"body"`
}

type chirpMsgError struct {
	Error string `json:"error"`
}

type chirpMessageValid struct {
	Valid        bool   `json:"valid"`
	Cleaned_body string `json:"cleaned_body"`
}

type userAdd struct {
	Email string `json:"email"`
}

func validate_chirp(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	msg := chirpMsg{}
	err := decoder.Decode(&msg)

	if err != nil {
		errorMsg := chirpMsgError{
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

	if len(msg.Body) > 140 {
		errorMsg := chirpMsgError{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(errorMsg)

		if err != nil {
			respondWithError(writer, http.StatusInternalServerError, "Error marshalling JSON during size check")
			return
		}

		respondWithJSON(writer, http.StatusBadRequest, data)
		return
	}

	msgCleaned := cleanBadWords(msg.Body)
	msgValid := chirpMessageValid{
		Valid:        true,
		Cleaned_body: msgCleaned,
	}
	data, err := json.Marshal(msgValid)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error marshalling JSON before sending response")
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
}

func (config *apiConfig) createUser(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	user := userAdd{}
	err := decoder.Decode(&user)

	if err != nil {
		errorMsg := chirpMsgError{
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
	createdUser, createError := config.db.CreateUser(request.Context(), createUser)
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

func (cfg *apiConfig) reset(writer http.ResponseWriter, request *http.Request) {
	if cfg.env == "dev" {
		deleteError := cfg.db.DeleteAllUsers(context.Background())

		if deleteError != nil {
			respondWithError(writer, http.StatusInternalServerError, "Unable to delete all users")
		}
	} else {
		respondWithError(writer, http.StatusForbidden, "Unable to perform this action in this environment")
	}
}
