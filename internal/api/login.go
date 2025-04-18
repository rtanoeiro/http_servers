package api

import (
	"encoding/json"
	"http_server/internal/auth"
	"http_server/internal/database"
	"log"
	"net/http"
	"time"
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
	log.Println("User details after login. \n- User:", userDetails.ID, "\n- Hashed Password:", userDetails.HashedPassword, "\n- Created At:", userDetails.CreatedAt, "\n- Updated At:", userDetails.UpdatedAt)

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
	log.Println("JWT Token Created with Success:", userJWTToken)

	userResfreshToken, errRefreshToken := MakeRefreshToken()
	if errRefreshToken != nil {
		respondWithError(writer, http.StatusUnauthorized, errRefreshToken.Error())
		return
	}
	log.Println("Refresh Token Created with Success:", userResfreshToken)
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     userResfreshToken,
		UserID:    userDetails.ID,
		ExpiresAt: time.Now().Add(time.Duration(60) * 60 * 24 * 60), //default 60 days expire
	}
	refreshTokenDetails, errRefreshToken := cfg.Db.CreateRefreshToken(request.Context(), refreshTokenParams)
	if errRefreshToken != nil {
		respondWithError(writer, http.StatusUnauthorized, errRefreshToken.Error())
		return
	}
	log.Println("Refresh Token Details Created with Success. \n- User:", refreshTokenDetails.UserID, "\n- Created At:", refreshTokenDetails.CreatedAt, "\n- Updated At:", refreshTokenDetails.UpdatedAt, "\n- Expires At:", refreshTokenDetails.ExpiresAt)

	loginResponse := UserResponse{
		ID:           userDetails.ID,
		CreatedAt:    userDetails.CreatedAt,
		UpdatedAt:    userDetails.UpdatedAt,
		Email:        user.Email,
		JWTToken:     &userJWTToken,
		RefreshToken: userResfreshToken,
	}

	loginBytes, marshalError := json.Marshal(loginResponse)
	if marshalError != nil {
		respondWithError(writer, http.StatusUnauthorized, marshalError.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, loginBytes)
}

func (cfg *ApiConfig) Refresh(writer http.ResponseWriter, request *http.Request) {
	refreshToken, errToken := GetBearerToken(request.Header)
	if errToken != nil {
		respondWithError(writer, http.StatusBadRequest, errToken.Error())
	}
	log.Println("Refresh Token from Request:", refreshToken)

	dbRefreshToken, dbError := cfg.Db.GetRefreshToken(request.Context(), refreshToken)
	if dbError != nil {
		respondWithError(writer, http.StatusUnauthorized, dbError.Error())
		return
	}
	log.Println("dbRefreshToken After Request:", dbRefreshToken)
	newAccessToken, errJWTToken := MakeJWT(dbRefreshToken.UserID, cfg.Secret)
	if errJWTToken != nil {
		respondWithError(writer, http.StatusUnauthorized, errJWTToken.Error())
		return
	}
	log.Println("New Access Token Created with Success:", newAccessToken)

	refreshResponse := RefreshResponse{
		Token: newAccessToken,
	}
	responseBytes, errMarshal := json.Marshal(refreshResponse)
	if errMarshal != nil {
		respondWithError(writer, http.StatusInternalServerError, errMarshal.Error())
	}
	respondWithJSON(writer, http.StatusOK, responseBytes)
}

func (cfg *ApiConfig) Revoke(writer http.ResponseWriter, request *http.Request) {
	refreshToken, errToken := GetBearerToken(request.Header)
	if errToken != nil {
		respondWithError(writer, http.StatusBadRequest, errToken.Error())
	}
	log.Println("Refresh Token from Request:", refreshToken)

	cfg.Db.RevokeRefreshToken(request.Context(), refreshToken)
	respondWithJSON(writer, http.StatusNoContent, []byte{})
}
