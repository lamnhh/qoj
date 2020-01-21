package src

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type UserClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func createToken(username string, jwtKey []byte, expirationTime time.Duration) string {
	claim := UserClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, _ := token.SignedString(jwtKey)
return tokenString
}

func CreateAccessToken(username string) string {
	return createToken(username, []byte(os.Getenv("JWT_ACCESS_SECRET")), 5 * time.Minute)
}

func CreateRefreshToken(username string) string {
	return createToken(username, []byte(os.Getenv("JWT_REFRESH_SECRET")), 7 * 24 * time.Hour)
}

func DecodeAccessToken(tokenString string) (string, error) {
	claim := UserClaim{}

	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	return claim.Username, nil
}

func DecodeRefreshToken(tokenString string) (string, error) {
	claim := UserClaim{}

	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	return claim.Username, nil
}