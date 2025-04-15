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

func ProcessChirp(request *http.Request) (int32, []byte, error) {

	decoder := json.NewDecoder(request.Body)
	msg := ChirpMsg{}
	err := decoder.Decode(&msg)

	if err != nil {
		errorMsg := ChirpMsgError{
			Error: "Something went wrong",
		}
		data, err := json.Marshal(errorMsg.Error)

		// Failed to marshal the error message
		if err != nil {
			return http.StatusInternalServerError, []byte{}, errors.New("Error marshalling JSON during initial request check")
		}
		// Error could be marshalled and is sent
		return http.StatusInternalServerError, data, nil
	}

	if len(msg.Body) > 140 {
		errorMsg := ChirpMsgError{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(errorMsg.Error)

		if err != nil {
			// Failed to marshal too long error msg
			return http.StatusInternalServerError, []byte{}, errors.New("Error marshalling JSON during size check")
		}
		return http.StatusBadRequest, data, nil
	}

	msgCleaned := CleanBadWords(msg.Body)
	msgValid := ChirpMessageValid{
		Valid:        true,
		Cleaned_body: msgCleaned,
	}
	data, err := json.Marshal(msgValid)
	if err != nil {
		return http.StatusInternalServerError, []byte{}, errors.New("Error marshalling JSON before sending response")
	}
	return http.StatusOK, data, nil
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
