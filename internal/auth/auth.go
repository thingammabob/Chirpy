package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeRefreshToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	output, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(1 * time.Hour)),
		Subject:   userID.String(),
	}).SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return output, nil
}
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, _ := token.Claims.(*jwt.RegisteredClaims)
	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}
	return userUUID, nil
}
func HashPassword(password string) (string, error) {
	hashed_password, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashed_password, nil
}

func CheckPaswordHash(password, hash string) (bool, error) {
	success, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return success, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	jwtArr := strings.Split(headers.Get("Authorization"), " ")
	if len(jwtArr) < 2 || jwtArr[0] != "Bearer" {
		return "", errors.New("Invalid authorization header")
	}
	jwt := jwtArr[1]
	if jwt == "" {
		return "", errors.New("JWT is empty")
	}
	return jwt, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	apiArr := strings.Split(headers.Get("Authorization"), " ")
	if len(apiArr) < 2 || apiArr[0] != "ApiKey" {
		return "", errors.New("Invalid authorization header")
	}
	api := apiArr[1]
	if api == "" {
		return "", errors.New("API key empty")
	}
	return api, nil
}
