package language

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getLanguage(ctx *gin.Context) {
	languageList, err := fetchAllLanguages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, languageList)
	}
}

func getLanguageId(ctx *gin.Context) {
	languageId, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid language ID"})
		return
	}

	language, err := FetchLanguageById(int(languageId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, language)
	}
}

func InitialiseRoutes(app *gin.RouterGroup) {
	app.GET("/language", getLanguage)
	app.GET("/language/:id", getLanguageId)
}