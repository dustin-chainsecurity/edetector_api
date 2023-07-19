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

type MissionResponse struct {
	IsSuccess   bool          `json:"isSuccess"`
	Message     string        `json:"message"`
	TaskID      []string      `json:"taskId"`
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

	tasks := []string{}
	for _, deviceId := range req.Devices {
		taskId, err := AddTask(deviceId, req.Action, messageMap[req.Action])
		if err != nil {
			Error.Handler(c, err, "Error adding task")
			return
		}
		tasks = append(tasks, taskId)
	}

	res := MissionResponse {
		IsSuccess: true,
		Message: "Success",
		TaskID: tasks,
	}
	c.JSON(http.StatusOK, res)
}