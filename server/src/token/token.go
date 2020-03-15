package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"qoj/server/config"
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

func ParseAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, err := parseUsernameFromToken(ctx)

		// If no token exists, or token is valid, ignore it
		if err != nil {
			ctx.Next()
			return
		}

		// Set username if possible
		ctx.Set("username", username)
		ctx.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, err := parseUsernameFromToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.Set("username", username)
		ctx.Next()
	}
}

// RequireAdmin is a middleware that verify if ctx.GetString(username) is an admin or not
// It is meant to be used AFTER RequireAuth()
func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.GetString("username")

		err := config.DB.QueryRow("SELECT username FROM users WHERE username = $1 AND is_admin = TRUE", username).
			Scan(&username)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Action requires admin privilege"})
			return
		}
		ctx.Next()
	}
}
