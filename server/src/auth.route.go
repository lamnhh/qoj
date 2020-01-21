package src

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

func postRegister(ctx *gin.Context) {
	var user User

	// If request body is not a JSON, return a 400.
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If some error occurs during user creation, return it back to user with a 500.
	if err := CreateNewUser(user); err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			// If err is an instance of pq.Error, it means the RAISE line in function create_user() has been called
			// Which, means that user.Username has been used before. In this case, return a 400.
			ctx.JSON(http.StatusBadRequest, gin.H{"error": pqErr.Hint})
		} else {
			// Otherwise, return a 500
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Set refresh token in cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "rftoken",
		Value:    CreateRefreshToken(user.Username),
		HttpOnly: true,
	})

	// Send access token back to user
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": CreateAccessToken(user.Username),
	})
}

func postLogin(ctx *gin.Context) {
	var userLogin UserLogin

	// If request body is not a JSON, return a 400.
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user with username `userLogin.Username`
	user, err := FindUserByUsername(userLogin.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			// If `username` does not exist, return a 404
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User \"" + userLogin.Username + "\" does not exist"})
		} else {
			// Otherwise, return a 500
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Set refresh token in cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "rftoken",
		Value:    CreateRefreshToken(user.Username),
		HttpOnly: true,
	})

	// Return (username, fullname)
	ctx.JSON(http.StatusOK, gin.H{
		"username":    user.Username,
		"fullname":    user.Fullname,
		"accessToken": CreateAccessToken(user.Username),
	})
}

func getRefresh(ctx *gin.Context) {
	// Get cookie `rftoken`
	cookie, err := ctx.Cookie("rftoken")
	if err != nil {
		// No cookie `rftoken`, return a 401
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	fmt.Println(cookie)

	// Decode the cookie above for `username`
	username, err := DecodeRefreshToken(cookie)
	if err != nil {
		// Cannot decode JWT, or JWT is expired, return 401
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set new refresh token in cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "rftoken",
		Value:    CreateRefreshToken(username),
		HttpOnly: true,
	})

	// Send access token back to user
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": CreateAccessToken(username),
	})
}

func getSecret(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"secret": "Hello there",
	})
}

func InitialiseAuthRoutes(app *gin.Engine) {
	app.POST("/api/register", postRegister)
	app.POST("/api/login", postLogin)
	app.GET("/api/refresh", getRefresh)

	// Protected route for token testing
	app.GET("/api/secret", AuthRequired(), getSecret)
}

