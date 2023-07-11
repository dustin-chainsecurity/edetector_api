package task

import (
	"edetector_API/internal/Error"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MissionRequest struct {
	Action      string      `json:"action"`
	Devices     []string    `json:"deviceId"`
}

var messageMap = map[string] string {
	"StartScan": "Ring0Process",
	"StartCollect": "null",
	"StartGetDrive": "null",
	"StartGetImage": "null",
	"Terminate": "null",
}

func SendMission(c *gin.Context) {

	var req MissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error.Handler(c, err, "Invalid request format")
		return
	}

	for _, deviceId := range req.Devices {
		err := AddTask(deviceId, req.Action, messageMap[req.Action])
		if err != nil {
			Error.Handler(c, err, "Error adding task")
			return
		}
	}

	res := TaskResponse {
		IsSuccess: true,
		Message: "Success",
	}
	c.JSON(http.StatusOK, res)
}