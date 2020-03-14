package auth

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"qoj/server/src/token"
	user2 "qoj/server/src/user"
)

func postLoginAdmin(ctx *gin.Context) {
	var userLogin user2.LoginAuth

	// If request body is not a JSON, return a 400.
	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := userLogin.Username
	password := userLogin.Password

	user, err := user2.FindAdminByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Username `%s` does not exist", username)})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong password"})
		return
	}

	// Set refresh token in cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "rftoken",
		Value:    token.CreateRefreshToken(username),
		HttpOnly: true,
	})

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": token.CreateAccessToken(username),
	})
}

func InitialiseAdminRoutes(app *gin.RouterGroup) {
	app.POST("/login", postLoginAdmin)
	app.GET("/refresh", getRefresh)
	app.GET("/logout", getLogout)
}