package testing

import (
	"edetector_API/api/task"
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateProgress(c *gin.Context) {

	// Receive request
	var req struct {
		TaskId   string `json:"taskId"`
		Progress int    `json:"progress"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	// Check taskId
	_, err := query.CheckDevice(req.TaskId)
	if err != nil {
		errhandler.Handler(c, err, "Error checking taskID")
		return
	}

	// Update progress
	err = query.UpdateTaskProgress(req.TaskId, req.Progress)
	if err != nil {
		errhandler.Handler(c, err, "Error updating task progress")
		return
	}

	// Send response
	res := task.TaskResponse{
		IsSuccess: true,
		Message:   "success",
	}
	c.JSON(http.StatusOK, res)
}
