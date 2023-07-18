package task

import (
	"edetector_API/internal/channel"
	"edetector_API/internal/Error"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Update(c *gin.Context) {

	// Receive request
	var req struct { DeviceId  string  `json:"deviceId"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		Error.Handler(c, err, "Invalid request format")
		return
	}

	// Send signal to websocket
	channel.SignalChannel <- []string{req.DeviceId}

	// Send response
	res := TaskResponse {
		IsSuccess: true,
		Message: "success",
	}
	c.JSON(http.StatusOK, res)
}
