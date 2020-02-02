package common

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParseQueryInt(ctx *gin.Context, name string, defaultValue int) int {
	val, err := strconv.ParseInt(ctx.Query(name), 10, 16)
	if err != nil {
		return defaultValue
	}
	return int(val)
}
