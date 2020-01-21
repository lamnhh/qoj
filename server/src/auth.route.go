package src

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
)

func postRegister(ctx *gin.Context) {
	var user User

	// If request body is not a JSON, return a 400.
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// If some error occurs during user creation, return it back to user with a 500.
	if err := CreateNewUser(user); err != nil {
		pqErr, ok := err.(*pq.Error)

		if ok {
			// If err is an instance of pq.Error, it means the RAISE line in function create_user() has been called
			// Which, means that user.Username has been used before. In this case, return a 400.
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": pqErr.Hint,
			})
		} else {
			// Otherwise, return a 500
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// TODO: create access token and refresh token
	// Access token is sent in response
	// Refresh token is stored in cookie
	ctx.Writer.WriteHeader(200)
}

func postLogin(ctx *gin.Context) {

}

func getRefresh(ctx *gin.Context) {

}

func InitialiseAuthRoute(app *gin.Engine) {
	app.POST("/api/register", postRegister)
	app.POST("/api/login", postLogin)
	app.GET("/api/refresh", getRefresh)
}
