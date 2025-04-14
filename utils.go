package main

import (
	"log"
	"net/http"
)

func respondWithError(writer http.ResponseWriter, code int, msg string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	log.Printf("Error marshalling JSON: %s", msg)
	writer.WriteHeader(code)
	writer.Write([]byte("Error marshalling JSON during initial request check"))
}

func respondWithJSON(writer http.ResponseWriter, code int, data []byte) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(code)
	writer.Write(data)
}
