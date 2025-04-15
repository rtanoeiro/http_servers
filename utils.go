package main

import (
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

func cleanBadWords(text string) string {

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
