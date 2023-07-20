package task

import (
	"edetector_API/internal/channel"
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/mariadb/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ScheduledCollect(c *gin.Context) {

	var req ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}

	for _, deviceId := range req.Devices {
		err := query.UpdateSchedule(deviceId, "collect_schedule", processSchedule(req.Date, req.Time))
		if err != nil {
			errhandler.Handler(c, err, "Error handling collect schedule")
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
