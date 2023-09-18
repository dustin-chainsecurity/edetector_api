package task

import (
	"edetector_API/internal/device"
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MissionRequest struct {
	Action  string   `json:"action"`
	Devices []string `json:"deviceId"`
	Type    string    `json:"type"`
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
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	// check devices
	if err := device.CheckAllID(req.Devices); err != nil {
		errhandler.Error(c, err, "Invalid device ID")
		return
	}
	// get userID
	userId, _ := c.Get("userID")
	uid, _ := userId.(int)
	// add tasks
	tasks := []string{}
	for _, deviceId := range req.Devices {
		taskId, err := addTask(uid, deviceId, req.Action, messageMap[req.Action], req.Type)
		if err != nil {
			errhandler.Error(c, err, "Error adding task")
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
