package api

import (
	"context"
	"encoding/json"
	"http_server/internal/auth"
	"http_server/internal/database"
	"log"
	"net/http"
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
	hashedPassword, hashError := auth.HashPassword(user.Password)

	if hashError != nil {
		respondWithError(writer, http.StatusInternalServerError, hashError.Error())
	}
	createUser := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          user.Email,
		HashedPassword: hashedPassword,
	}
	createdUser, createError := config.Db.CreateUser(request.Context(), createUser)
	if createError != nil {
		log.Println("CreateUser error:", createError)
		respondWithError(writer, http.StatusInternalServerError, "Unable to create user")
		return
	}

	returnUser := UserResponse{
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

func (cfg *ApiConfig) ResetUsers(writer http.ResponseWriter, request *http.Request) {
	if cfg.Env == "dev" {
		deleteError := cfg.Db.DeleteAllUsers(context.Background())

		if deleteError != nil {
			respondWithError(writer, http.StatusInternalServerError, deleteError.Error())
		}
	} else {
		respondWithError(writer, http.StatusForbidden, "Unable to perform this action in this environment")
	}
	respondWithJSON(writer, http.StatusOK, []byte{})
}
