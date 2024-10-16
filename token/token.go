package token

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pratikdev/url-shortner-with-go/customErrors"
)

type claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type tokenResponse struct {
	Value          string
	ExpirationTime time.Time
}

func GetToken(username string) (tokenResponse, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		err = &customErrors.CustomError{Code: http.StatusInternalServerError, Message: err.Error()}
		return tokenResponse{}, err
	}

	return tokenResponse{
		Value:          tokenString,
		ExpirationTime: expirationTime,
	}, nil
}

func ValidateToken(token string) (bool, error) {
	tokenStr := token
	claims := &claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !tkn.Valid {
		return false, &customErrors.CustomError{Code: http.StatusBadRequest, Message: "Unauthorized"}
	}

	return tkn.Valid, nil
}
