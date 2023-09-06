package analysis

import (
	"edetector_API/internal/device"
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

type detectDetailResponse struct {
	IsSuccess bool            `json:"isSuccess"`
	Data      []device.Detail `json:"data"`
}

func DeviceDetail(c *gin.Context) {

	// Process device data
	raw_devices, err := query.LoadAllDeviceInfo()
	if err != nil {
		errhandler.Error(c, err, "Error loading raw device data")
	}

	devices := []device.Detail{}
	for _, r := range raw_devices {
		d, err := device.ProcessDeviceDetail(r)
		if err != nil {
			errhandler.Error(c, err, "Error processing device data")
			return
		}
		devices = append(devices, d)
	}

	// Create the response object
	res := detectDetailResponse{
		IsSuccess: true,
		Data:      devices,
	}

	c.JSON(http.StatusOK, res)
}
