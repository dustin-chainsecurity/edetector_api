package searchEvidence

import (
	"edetector_API/internal/device"
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type detectDevicesResponse struct {
	IsSuccess bool            `json:"isSuccess"`
	Data      []device.Device `json:"data"`
}

func DetectDevices(c *gin.Context) {

	// Process device data
	raw_devices, err := query.LoadAllDeviceInfo()
	if err != nil {
		errhandler.Handler(c, err, "Error loading raw device data")
	}

	devices := []device.Device{}
	for _, r := range raw_devices {
		d, err := device.ProcessRawDevice(r)
		if err != nil {
			errhandler.Handler(c, err, "Error processing device data")
			return
		}
		devices = append(devices, d)
	}

	// Create the response object
	res := detectDevicesResponse{
		IsSuccess: true,
		Data:      devices,
	}
	logger.Debug("DetectDevices response: ", zap.Any("res", res))

	c.JSON(http.StatusOK, res)
}
