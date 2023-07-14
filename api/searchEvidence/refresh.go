package searchEvidence

import (
	"edetector_API/internal/Error"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

type refreshResponse struct {
	IsSuccess    bool       `json:"isSuccess"`
	Data         device     `json:"data"`
}

func Refresh(c *gin.Context) {

	// Receive request
	var req struct { DeviceId  string  `json:"deviceId"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		Error.Handler(c, err, "Invalid request format")
		return
	}

	raw_device, err := query.LoadDeviceInfo(req.DeviceId)
	if err != nil {
		Error.Handler(c, err, "Error loading raw device data")
	}
	device, err := processRawDevice(raw_device)
	if err != nil {
		Error.Handler(c, err, "Error processing device data")
	}

	// Send response
	res := refreshResponse {
		IsSuccess: true,
		Data: device,
	}
	c.JSON(http.StatusOK, res)
}