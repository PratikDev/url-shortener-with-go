package token

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pratikdev/url-shortner-with-go/customErrors"
)

type customClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type tokenResponse struct {
	Value          string
	ExpirationTime time.Time
}

// get token
func GetToken(userId string) (tokenResponse, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &customClaims{
		UserId: userId,
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

type validationResponse struct {
	IsValid bool   `json:"isValid"`
	UserID  string `json:"userId"`
}

// validate token
func ValidateToken(tokenStr string) (validationResponse, error) {
	claims := &customClaims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !tkn.Valid {
		return validationResponse{}, &customErrors.CustomError{Code: http.StatusBadRequest, Message: "Unauthorized"}
	}

	return validationResponse{
		IsValid: tkn.Valid,
		UserID:  claims.UserId,
	}, nil
}

// refresh token when only 5 minutes are left
func RefreshToken(tokenStr string) (tokenResponse, error) {
	claims := &customClaims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !tkn.Valid {
		return tokenResponse{}, &customErrors.CustomError{Code: http.StatusBadRequest, Message: "Unauthorized"}
	}

	if time.Until(claims.ExpiresAt.Time).Minutes() <= 5 {
		return GetToken(claims.UserId)
	}

	return tokenResponse{
		Value:          tokenStr,
		ExpirationTime: claims.ExpiresAt.Time,
	}, nil
}
