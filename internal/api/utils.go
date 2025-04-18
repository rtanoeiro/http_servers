package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	claim := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(60) * time.Minute)), //default 60 min expire
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

func GetAuthorizationField(headers http.Header) (string, error) {

	fullHeader := headers.Get("Authorization")
	headerFields := strings.Fields(fullHeader)
	if len(headerFields) < 2 {
		return "", errors.New("invalid header format, unable to perform action")
	}

	token := headerFields[1]
	return token, nil
}

func MakeRefreshToken() (string, error) {
	tokenBytes := make([]byte, 256)
	number, _ := rand.Read(tokenBytes)
	fmt.Println(number)
	myTokenString := hex.EncodeToString(tokenBytes)

	return myTokenString, nil
}
