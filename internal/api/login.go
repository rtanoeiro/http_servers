package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"http_server/internal/auth"
	"net/http"
	"strings"
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

	userJWTToken, errJWTToken := MakeJWT(userDetails.ID, cfg.Secret)
	if errJWTToken != nil {
		respondWithError(writer, http.StatusUnauthorized, errJWTToken.Error())
		return
	}

	userResfreshToken, errRefreshToken := MakeRefreshToken()
	if errRefreshToken != nil {
		respondWithError(writer, http.StatusUnauthorized, errRefreshToken.Error())
		return
	}
	loginResponse := UserResponse{
		ID:           userDetails.ID,
		CreatedAt:    userDetails.CreatedAt,
		UpdatedAt:    userDetails.UpdatedAt,
		Email:        user.Email,
		Token:        &userJWTToken,
		RefreshToken: userResfreshToken,
	}

	loginBytes, marshalError := json.Marshal(loginResponse)
	if marshalError != nil {
		respondWithError(writer, http.StatusUnauthorized, marshalError.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, loginBytes)
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

func GetBearerToken(headers http.Header) (string, error) {

	fullHeader := headers.Get("Authorization")
	headerFields := strings.Fields(fullHeader)
	if len(headerFields) < 2 {
		return "", nil
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
