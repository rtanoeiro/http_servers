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

	httpStatusCode, chirpRequest, valError := ProcessChirp(request)

	if valError != nil {
		respondWithError(writer, int(httpStatusCode), valError.Error())
	}
	args := database.InsertChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      chirpRequest.Body,
		UserID:    chirpRequest.UserID,
	}
	chirp, errorInsert := cfg.Db.InsertChirp(request.Context(), args)

	if errorInsert != nil {
		respondWithError(writer, http.StatusInternalServerError, errorInsert.Error())
	}
	chirpResponse := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	chirpBytes, marshalError := json.Marshal(chirpResponse)

	if marshalError != nil {
		respondWithError(writer, http.StatusInternalServerError, marshalError.Error())
	}
	respondWithJSON(writer, http.StatusCreated, chirpBytes)
}

func (cfg *ApiConfig) GetAllChirps(writer http.ResponseWriter, request *http.Request) {
	allChirps, chirpError := cfg.Db.GetAllChirps(request.Context())

	if chirpError != nil {
		respondWithError(writer, http.StatusInternalServerError, chirpError.Error())
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
	}
	respondWithJSON(writer, http.StatusOK, chirpsBytes)
}

func (cfg *ApiConfig) GetSingleChirp(writer http.ResponseWriter, request *http.Request) {
	chirpID := uuid.MustParse(request.PathValue("chirpID"))
	fmt.Println("Parsing ChirpID:", chirpID)
	chirp, chirpError := cfg.Db.GetSingleChirp(request.Context(), chirpID)

	if chirpError != nil {
		respondWithError(writer, http.StatusNotFound, chirpError.Error())
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
	}
	respondWithJSON(writer, http.StatusOK, chirpsBytes)
}
