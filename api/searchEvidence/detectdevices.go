package searchEvidence

import (
	"edetector_API/internal/device"
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
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

	c.JSON(http.StatusOK, res)
}
