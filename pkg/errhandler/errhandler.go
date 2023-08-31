package errhandler

import (
	"edetector_API/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context, err error, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": false,
		"message":   msg + ": " + err.Error(),
	})
	logger.Error(msg + ": " + err.Error())
}