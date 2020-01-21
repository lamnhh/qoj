package src

import (
	"database/sql"
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

	// TODO: create access token and refresh token
	// Access token is sent in response
	// Refresh token is stored in cookie
	ctx.Writer.WriteHeader(200)
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

	// Return (username, fullname)
	ctx.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"fullname": user.Fullname,
	})
}

func getRefresh(ctx *gin.Context) {

}

func InitialiseAuthRoute(app *gin.Engine) {
	app.POST("/api/register", postRegister)
	app.POST("/api/login", postLogin)
	app.GET("/api/refresh", getRefresh)
}
