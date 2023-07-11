package Error

import (
	"edetector_API/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Error error

func Handler(c *gin.Context, err error, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{msg: err.Error()})
	logger.Error(msg + ": " + err.Error())
}