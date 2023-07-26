package task

import (
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MissionRequest struct {
	Action  string   `json:"action"`
	Devices []string `json:"deviceId"`
}

type MissionResponse struct {
	IsSuccess bool     `json:"isSuccess"`
	Message   string   `json:"message"`
	TaskID    []string `json:"taskId"`
}

var messageMap = map[string]string{
	"StartScan":     "Ring0Process",
	"StartCollect":  "null",
	"StartGetDrive": "null",
	"StartGetImage": "null",
	"Terminate":     "null",
}

func SendMission(c *gin.Context) {

	var req MissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	// check devices
	err := query.CheckAllDevice(req.Devices)
	if err != nil {
		errhandler.Handler(c, err, "Error checking deviceID")
		return
	}

	tasks := []string{}
	for _, deviceId := range req.Devices {
		taskId, err := addTask(deviceId, req.Action, messageMap[req.Action])
		if err != nil {
			errhandler.Handler(c, err, "Error adding task")
			return
		}
		tasks = append(tasks, taskId)
	}

	res := MissionResponse{
		IsSuccess: true,
		Message:   "Success",
		TaskID:    tasks,
	}
	c.JSON(http.StatusOK, res)
}
