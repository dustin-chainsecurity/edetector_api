package task

import (
	"edetector_API/internal/channel"
	"edetector_API/internal/device"
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateTask(c *gin.Context) {
	// Receive request
	var req struct {
		DeviceId string `json:"deviceId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	// Check deviceId
	if err := device.CheckID(req.DeviceId); err != nil {
		errhandler.Handler(c, err, "Invalid device ID")
		return
	}
	// Send signal to websocket
	channel.SignalChannel <- []string{req.DeviceId}
	// Send response
	res := TaskResponse{
		IsSuccess: true,
		Message:   "success",
	}
	c.JSON(http.StatusOK, res)
}