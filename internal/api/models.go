package api

import (
	"http_server/internal/database"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Env            string
	Secret         string
}

type ChirpMsgError struct {
	Error string `json:"error"`
}

type ChirpRequest struct {
	Body string `json:"body"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

type RefreshResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	JWTToken     *string   `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}
type UserAdd struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UpdateUser struct {
	Email string `json:"email"`
}

type PolkaWebHook struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}
