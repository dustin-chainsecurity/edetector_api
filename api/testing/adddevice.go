package testing

import (
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deviceInfo struct {
	DeviceID   string `json:"deviceID"`
	DeviceName string `json:"deviceName"`
	IP         string `json:"ip"`
}

type addDeviceRequest struct {
	Devices   []deviceInfo  `json:"devices"`
}

func AddDevice(c *gin.Context) {
	var req addDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Error binding JSON")
		return
	}
	for _, d := range req.Devices {
		err := query.AddDevice(d.DeviceID, d.DeviceName, d.IP)
		if err != nil {
			errhandler.Handler(c, err, "Error adding device")
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}