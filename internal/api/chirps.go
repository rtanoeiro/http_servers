package api

import (
	"encoding/json"
	"fmt"
	"http_server/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) InsertChirp(writer http.ResponseWriter, request *http.Request) {

	token, _ := GetBearerToken(request.Header)
	jwtUserID, errorJWT := ValidateJWT(token, cfg.Secret)

	if errorJWT != nil {
		respondWithError(writer, http.StatusUnauthorized, errorJWT.Error())
		return
	}

	httpStatusCode, chirpRequest, valError := ProcessChirp(request)
	fmt.Println("Chirp Procesed: \n - Body:", chirpRequest.Body, "\n - Error:", valError)
	fmt.Println("User ID From Token: ", jwtUserID, "\nError from JWT Token", errorJWT)

	if valError != nil {
		respondWithError(writer, int(httpStatusCode), valError.Error())
		return
	}
	args := database.InsertChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      chirpRequest.Body,
		UserID:    jwtUserID,
	}
	chirp, errorInsert := cfg.Db.InsertChirp(request.Context(), args)

	if errorInsert != nil {
		respondWithError(writer, http.StatusInternalServerError, errorInsert.Error())
		return
	}
	chirpResponse := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    jwtUserID,
	}
	chirpBytes, marshalError := json.Marshal(chirpResponse)

	if marshalError != nil {
		respondWithError(writer, http.StatusInternalServerError, marshalError.Error())
		return
	}
	respondWithJSON(writer, http.StatusCreated, chirpBytes)
}

func (cfg *ApiConfig) GetAllChirps(writer http.ResponseWriter, request *http.Request) {
	allChirps, chirpError := cfg.Db.GetAllChirps(request.Context())

	if chirpError != nil {
		respondWithError(writer, http.StatusInternalServerError, chirpError.Error())
		return
	}

	chirpsResponse := make([]ChirpResponse, len(allChirps))
	for i, chirp := range allChirps {
		chirpsResponse[i] = ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}
	chirpsBytes, err := json.Marshal(chirpsResponse)

	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, chirpsBytes)
}

func (cfg *ApiConfig) GetSingleChirp(writer http.ResponseWriter, request *http.Request) {
	chirpID := uuid.MustParse(request.PathValue("chirpID"))
	chirp, chirpError := cfg.Db.GetSingleChirp(request.Context(), chirpID)

	if chirpError != nil {
		respondWithError(writer, http.StatusNotFound, chirpError.Error())
		return
	}

	chirpResponse := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	chirpsBytes, err := json.Marshal(chirpResponse)

	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, chirpsBytes)
}
