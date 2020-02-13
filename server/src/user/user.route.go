package user

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qoj/server/config"
	"qoj/server/src/problem"
	"qoj/server/src/token"
)

func getUser(ctx *gin.Context) {
	username := ctx.GetString("username")
	user, err := FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("User `%s` does not exist", username),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		user.Password = ""
		ctx.JSON(http.StatusOK, user)
	}
}

func getUserPublic(ctx *gin.Context) {
	username := ctx.Param("username")
	user, err := FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("User `%s` does not exist", username),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	} else {
		user.Password = ""
		ctx.JSON(http.StatusOK, user)
	}
}

func getUserSolved(ctx *gin.Context) {
	username := ctx.Param("username")
	_, err := FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("User `%s` does not exist", username),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	problemList, err := problem.FetchSolvedProblemsOfUser(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, problemList)
	}
}

func getUserPartial(ctx *gin.Context) {
	username := ctx.Param("username")
	_, err := FindUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("User `%s` does not exist", username),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	problemList, err := problem.FetchPartiallySolvedProblemsOfUser(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, problemList)
	}
}

func postUserProfilePicture(ctx *gin.Context) {
	username := ctx.GetString("username")
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if path, err := updateProfilePicture(username, file); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"path": path})
	}
}

func patchUser(ctx *gin.Context) {
	username := ctx.GetString("username")

	patch := PatchUser{}
	if err := ctx.ShouldBindJSON(&patch); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	modifiers := make(map[string]interface{})
	if val := patch.Fullname; val != "" {
		modifiers["fullname"] = val
	}

	user, err := UpdateUserProfile(username, modifiers)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		user.Password = ""
		ctx.JSON(http.StatusOK, user)
	}
}

func putUserPassword(ctx *gin.Context) {
	username := ctx.GetString("username")

	body := PutPasswordUser{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, code, err := Login(username, body.OldPassword)
	if err != nil {
		ctx.JSON(code, gin.H{"error": err.Error()})
		return
	}

	hashedPassword := hashPassword(body.NewPassword)
	_, err = config.DB.Exec("UPDATE users SET password = $1 WHERE username = $2", hashedPassword, username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func InitialiseUserRoutes(app *gin.Engine) {
	initialiseAvatarLocks()

	app.GET("/api/user", token.RequireAuth(), getUser)
	app.GET("/api/user/:username/public", getUserPublic)
	app.GET("/api/user/:username/solved", getUserSolved)
	app.GET("/api/user/:username/partial", getUserPartial)

	// Endpoint to upload profile picture
	app.POST("/api/user/profile-picture", token.RequireAuth(), postUserProfilePicture)

	// Update current user's information
	app.PATCH("/api/user", token.RequireAuth(), patchUser)

	// Update password
	app.PUT("/api/user/password", token.RequireAuth(), putUserPassword)
}