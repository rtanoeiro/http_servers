package api

import (
	"encoding/json"
	"http_server/internal/auth"
	"net/http"
)

func (cfg *ApiConfig) Login(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	user := UserLogin{}
	err := decoder.Decode(&user)

	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	// hashProvided, errProvided := auth.HashPassword(user.Password)
	userDetails, loginErr := cfg.Db.CheckUserLogin(request.Context(), user.Email)
	if loginErr != nil {
		respondWithError(writer, http.StatusUnauthorized, loginErr.Error())
		return

	}

	results := auth.CheckPasswordHash(user.Password, userDetails.HashedPassword)
	if results != nil {
		respondWithError(writer, http.StatusUnauthorized, results.Error())
		return
	}

	loginResponse := UserResponse{
		ID:        userDetails.ID,
		CreatedAt: userDetails.CreatedAt,
		UpdatedAt: userDetails.UpdatedAt,
		Email:     user.Email,
	}
	loginBytes, marshalError := json.Marshal(loginResponse)

	if marshalError != nil {
		respondWithError(writer, http.StatusUnauthorized, marshalError.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, loginBytes)
}
