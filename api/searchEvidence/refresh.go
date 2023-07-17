package searchEvidence

import (
	"edetector_API/internal/Error"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshResponse struct {
	IsSuccess    bool       `json:"isSuccess"`
	Data         []device     `json:"data"`
}

func Refresh(c *gin.Context) {

	// Receive request
	var req struct { Devices  []string  `json:"deviceId"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		Error.Handler(c, err, "Invalid request format")
		return
	}

	var devices = []device{}
	for _, deviceId := range req.Devices {
		raw_device, err := query.LoadDeviceInfo(deviceId)
		if err != nil {
			Error.Handler(c, err, "Error loading raw device data")
			return
		}
		device, err := processRawDevice(raw_device)
		if err != nil {
			Error.Handler(c, err, "Error processing device data")
			return
		}
		devices = append(devices, device)
	}

	// Send response
	res := refreshResponse {
		IsSuccess: true,
		Data: devices,
	}
	c.JSON(http.StatusOK, res)
}