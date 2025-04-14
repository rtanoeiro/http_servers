package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type chirpMsg struct {
	Body string `json:"body"`
}

type chirpMsgError struct {
	Error string `json:"error"`
}

type chirpMessageValid struct {
	Valid bool `json:"valid"`
}

func validate_chirp(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	msg := chirpMsg{}
	err := decoder.Decode(&msg)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		errorMsg := chirpMsgError{
			Error: "Something went wrong",
		}
		data, err := json.Marshal(errorMsg)

		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Error marshalling JSON during initial request check"))
			return
		}
		writer.Write(data)
		return
	}

	if len(msg.Body) > 140 {
		writer.WriteHeader(http.StatusBadRequest)
		errorMsg := chirpMsgError{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(errorMsg)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Error marshalling JSON during size check"))
			return
		}

		writer.Write(data)
		return
	}

	msgValid := chirpMessageValid{
		Valid: true,
	}
	writer.WriteHeader(http.StatusOK)
	data, err := json.Marshal(msgValid)

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("Error marshalling JSON before sending data to client"))
		return
	}
	writer.Write(data)
}
