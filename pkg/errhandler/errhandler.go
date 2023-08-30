package errhandler

import (
	"edetector_API/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context, err error, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{msg: err.Error()})
	logger.Error(msg + ": " + err.Error())
}