package task

import (
	"edetector_API/internal/Error"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ScheduledScan(c *gin.Context) {

	var req ScheduleScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error.Handler(c, err, "Invalid request format")
		return
	}

	for _, deviceId := range req.Devices {
		err := query.UpdateSchedule(deviceId, "scan_schedule", req.Time)
		if err != nil {
			Error.Handler(c, err, "Error handling scan schedule")
			return
		}
	}

	res := TaskResponse {
		IsSuccess: true,
		Message: "success",
	}
	c.JSON(http.StatusOK, res)
}