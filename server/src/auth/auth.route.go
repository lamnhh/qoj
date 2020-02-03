package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"qoj/server/src/token"
	user2 "qoj/server/src/user"
	"time"
)

func postRegister(ctx *gin.Context) {
	var user user2.User

	// If request body is not a JSON, return a 400.
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If some error occurs during user creation, return it back to user with a 500.
	if err := user2.CreateNewUser(user); err != nil {
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
		Value:    token.CreateRefreshToken(user.Username),
		HttpOnly: true,
	})

	// Send access token back to user
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": token.CreateAccessToken(user.Username),
	})
}

func postLogin(ctx *gin.Context) {
	var userLogin user2.LoginAuth

	// If request body is not a JSON, return a 400.
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user with username `userLogin.Username`
	user, code, err := user2.Login(userLogin.Username, userLogin.Password)
	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	// Set refresh token in cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "rftoken",
		Value:    token.CreateRefreshToken(user.Username),
		HttpOnly: true,
	})

	// Return (username, fullname)
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": token.CreateAccessToken(user.Username),
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

	// Decode the cookie above for `username`
	username, err := token.DecodeRefreshToken(cookie)
	if err != nil {
		// Cannot decode JWT, or JWT is expired, return 401
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set new refresh token in cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "rftoken",
		Value:    token.CreateRefreshToken(username),
		HttpOnly: true,
	})

	// Send access token back to user
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": token.CreateAccessToken(username),
	})
}

func getLogout(ctx *gin.Context) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "rftoken",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
	})
	ctx.JSON(http.StatusOK, gin.H{})
}

func InitialiseAuthRoutes(app *gin.Engine) {
	app.POST("/api/register", postRegister)
	app.POST("/api/login", postLogin)
	app.GET("/api/refresh", getRefresh)
	app.GET("/api/logout", getLogout)
}
