package searchEvidence

import (
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshResponse struct {
	IsSuccess bool     `json:"isSuccess"`
	Data      []device `json:"data"`
}

func Refresh(c *gin.Context) {

	// Receive request
	var req struct {
		Devices []string `json:"deviceId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	fmt.Println("Request content: ", req)

	var devices = []device{}
	for _, deviceId := range req.Devices {
		raw_device, err := query.LoadDeviceInfo(deviceId)
		if err != nil {
			errhandler.Handler(c, err, "Error loading raw device data")
			return
		}
		device, err := processRawDevice(raw_device)
		if err != nil {
			errhandler.Handler(c, err, "Error processing device data")
			return
		}
		devices = append(devices, device)
	}

	// Send response
	res := refreshResponse{
		IsSuccess: true,
		Data:      devices,
	}
	c.JSON(http.StatusOK, res)
}
