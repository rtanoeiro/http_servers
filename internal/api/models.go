package api

import (
	"http_server/internal/database"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Env            string
}

type ChirpMsg struct {
	Body string `json:"body"`
}

type ChirpMsgError struct {
	Error string `json:"error"`
}

type ChirpMessageValid struct {
	Valid        bool   `json:"valid"`
	Cleaned_body string `json:"cleaned_body"`
}

type UserAdd struct {
	Email string `json:"email"`
}
