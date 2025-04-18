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

func ProcessChirp(writer http.ResponseWriter, request *http.Request) (ChirpRequest, error) {

	decoder := json.NewDecoder(request.Body)
	chirpRequest := ChirpRequest{}
	err := decoder.Decode(&chirpRequest)

	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
		return ChirpRequest{Body: ""}, errors.New(err.Error())
	}

	if len(chirpRequest.Body) > 140 {
		respondWithError(writer, http.StatusBadRequest, "chirp is too long")
		return ChirpRequest{Body: ""}, errors.New("chirp is too long")
	}
	msgCleaned := CleanBadWords(chirpRequest.Body)
	cleanChirp := ChirpRequest{Body: msgCleaned}
	return cleanChirp, nil
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
