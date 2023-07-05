package task

import (
	"edetector_API/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartCollect(c *gin.Context) {

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	response, err := SendTCPRequest(c, "StartCollect", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		logger.Error("Error sending TCP request: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}