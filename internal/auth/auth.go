package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	output, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
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
	if len(jwtArr) != 2 {
		return "", errors.New("Invalid authorization header")
	}
	jwt := jwtArr[1]
	if jwt == "" {
		return "", errors.New("Invalid jwt")
	}
	return jwt, nil
}
