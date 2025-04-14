package main

import (
	"encoding/json"
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

	if err != nil {
		errorMsg := chirpMsgError{
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
		errorMsg := chirpMsgError{
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

	msgValid := chirpMessageValid{
		Valid: true,
	}
	data, err := json.Marshal(msgValid)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, "Error marshalling JSON before sending response")
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
}
