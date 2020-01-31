package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
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

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
			return
		}

		// Split authHeader into 2 parts: "Bearer" and the actual token
		authHeaderToken := strings.Split(authHeader, " ")
		if len(authHeaderToken) != 2 {
			// If there aren't exactly 2 parts, return a 401
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Extract access token as the 2nd element from authHeaderToken
		accessToken := authHeaderToken[1]
		username, err := DecodeAccessToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Save `username` into context for later use
		ctx.Set("username", username)
		ctx.Next()
	}
}
