package api

import (
	"encoding/json"
	"http_server/internal/auth"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) Login(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	user := UserLogin{}
	err := decoder.Decode(&user)

	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	// hashProvided, errProvided := auth.HashPassword(user.Password)
	userDetails, loginErr := cfg.Db.CheckUserLogin(request.Context(), user.Email)
	if loginErr != nil {
		respondWithError(writer, http.StatusUnauthorized, loginErr.Error())
		return

	}

	results := auth.CheckPasswordHash(user.Password, userDetails.HashedPassword)
	if results != nil {
		respondWithError(writer, http.StatusUnauthorized, results.Error())
		return
	}

	loginResponse := UserResponse{
		ID:        userDetails.ID,
		CreatedAt: userDetails.CreatedAt,
		UpdatedAt: userDetails.UpdatedAt,
		Email:     user.Email,
	}
	loginBytes, marshalError := json.Marshal(loginResponse)

	if marshalError != nil {
		respondWithError(writer, http.StatusUnauthorized, marshalError.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, loginBytes)
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claim := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := jwtToken.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Create a new token object by parsing the token string
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for validation
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	subject, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.Nil, err
	}
	return uuid.MustParse(subject), nil
}
