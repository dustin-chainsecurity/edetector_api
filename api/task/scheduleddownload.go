package task

import (
	"edetector_API/internal/channel"
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ScheduledDownload(c *gin.Context) {

	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	fmt.Println("Request content: ", req)

	for _, deviceId := range req.Devices {
		err := query.UpdateSchedule(deviceId, "file_schedule", processSchedule(req.Date, req.Time))
		if err != nil {
			errhandler.Handler(c, err, "Error handling file schedule")
			return
		}
	}

	channel.SignalChannel <- req.Devices

	res := TaskResponse{
		IsSuccess: true,
		Message:   "success",
	}
	c.JSON(http.StatusOK, res)
}
