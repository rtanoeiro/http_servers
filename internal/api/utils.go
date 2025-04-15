package api

import (
	"encoding/json"
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

func ValidateChirp(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	msg := ChirpMsg{}
	err := decoder.Decode(&msg)

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

	if len(msg.Body) > 140 {
		errorMsg := ChirpMsgError{
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

	msgCleaned := CleanBadWords(msg.Body)
	msgValid := ChirpMessageValid{
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
