package searchEvidence

import (
	"edetector_API/internal/device"
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshResponse struct {
	IsSuccess bool            `json:"isSuccess"`
	Data      []device.Device `json:"data"`
}

func Refresh(c *gin.Context) {

	// receive request
	var req struct {
		Devices []string `json:"deviceId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	// check devices
	if err := device.CheckAllID(req.Devices); err != nil {
		errhandler.Handler(c, err, "Invalid device ID")
		return
	}

	// load device data
	var devices = []device.Device{}
	for _, deviceId := range req.Devices {
		raw_device, err := query.LoadDeviceInfo(deviceId)
		if err != nil {
			errhandler.Handler(c, err, "Error loading raw device data")
			return
		}
		d, err := device.ProcessRawDevice(raw_device)
		if err != nil {
			errhandler.Handler(c, err, "Error processing device data")
			return
		}
		devices = append(devices, d)
	}

	// send response
	res := refreshResponse{
		IsSuccess: true,
		Data:      devices,
	}
	c.JSON(http.StatusOK, res)
}
