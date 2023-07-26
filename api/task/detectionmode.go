package task

import (
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DetectionModeRequest struct {
	Mode    bool     `json:"mode"`
	Devices []string `json:"deviceId"`
}

func DetectionMode(c *gin.Context) {

	var req DetectionModeRequest
	var message = "Success"
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	// check devices
	err := query.CheckAllDevice(req.Devices)
	if err != nil {
		errhandler.Handler(c, err, "Error checking deviceID")
		return
	}

	for _, deviceId := range req.Devices {

		// fetch agent current setting
		var m int
		if req.Mode {
			m = 1
		} else {
			m = 0
		}
		process, network, err := query.GetDetectMode(deviceId)
		if err != nil {
			errhandler.Handler(c, err, "Error retrieving current client setting")
			return
		}

		// check if the changes are needed
		if m != process || m != network {
			_, err := addTask(deviceId, "ChangeDetectMode", fmt.Sprintf("%d|%d", m, m))
			if err != nil {
				errhandler.Handler(c, err, "Error adding ChangeDetectMode task")
				return
			}
			err = query.UpdateDetectMode(m, deviceId)
			if err != nil {
				errhandler.Handler(c, err, "Error updating client_setting table")
				return
			}
		} else {
			message = "No changes needed"
		}
	}

	res := TaskResponse{
		IsSuccess: true,
		Message:   message,
	}
	c.JSON(http.StatusOK, res)
}