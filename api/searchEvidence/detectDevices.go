package searchEvidence

import (
	"edetector_API/internal/Error"
	"edetector_API/internal/packet"
	"net/http"
	"github.com/gin-gonic/gin"
)

type detectDevicesResponse struct {
	IsSuccess    bool              `json:"isSuccess"`
	Data         []packet.Device   `json:"data"`
}

func DetectDevices(c *gin.Context) {

	// Process device data
	devices, err := packet.ProcessDeviceData()
	if err != nil {
		Error.Handler(c, err, "Error processing device data")
		return
	}

	// Create the response object
	res := detectDevicesResponse{
		IsSuccess:    true,
		Data:         devices,
	}

	c.JSON(http.StatusOK, res)
}