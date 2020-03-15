package token

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func parseUsernameFromToken(ctx *gin.Context) (string, error) {
	// Extract header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("Access token required")
	}

	// Split authHeader into 2 parts: "Bearer" and the actual token
	authHeaderToken := strings.Split(authHeader, " ")
	if len(authHeaderToken) != 2 {
		// If there aren't exactly 2 parts, return a 401
		return "", errors.New("Invalid token")
	}

	// Decode token to get username
	accessToken := authHeaderToken[1]
	username, err := DecodeAccessToken(accessToken)
	if err != nil {
		return "", err
	}

	return username, nil
}
