package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func respondWithError(writer http.ResponseWriter, code int, msg string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(code)
	writer.Write([]byte(msg))
}

func respondWithJSON(writer http.ResponseWriter, code int, data []byte) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(code)
	writer.Write(data)
}

func ProcessChirp(request *http.Request) (int32, ChirpRequest, error) {

	decoder := json.NewDecoder(request.Body)
	chirpRequest := ChirpRequest{}
	err := decoder.Decode(&chirpRequest)

	if err != nil {
		errorMsg := ChirpMsgError{
			Error: "Something went wrong",
		}
		return http.StatusInternalServerError, ChirpRequest{
			Body:   "",
			UserID: chirpRequest.UserID,
		}, errors.New(errorMsg.Error)
	}

	if len(chirpRequest.Body) > 140 {
		errorMsg := ChirpMsgError{
			Error: "Chirp is too long",
		}
		return http.StatusBadRequest, ChirpRequest{
			Body:   "",
			UserID: chirpRequest.UserID}, errors.New(errorMsg.Error)
	}

	msgCleaned := CleanBadWords(chirpRequest.Body)
	cleanChirp := ChirpRequest{
		Body:   msgCleaned,
		UserID: chirpRequest.UserID,
	}
	return http.StatusOK, cleanChirp, nil
}

func CleanBadWords(text string) string {

	if len(text) == 0 {
		return text
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	splitWords := strings.Fields(text)

	for i, word := range splitWords {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				splitWords[i] = "****"
			}
		}
	}
	finalWords := strings.Join(splitWords, " ")
	fmt.Println(finalWords)
	return finalWords
}
