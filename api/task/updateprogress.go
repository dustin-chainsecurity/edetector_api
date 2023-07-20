package task

import (
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/mariadb"
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

	query := "UPDATE task SET progress = ? WHERE task_id = ?"
	_, err := mariadb.DB.Exec(query, req.Progress, req.TaskId)
	if err != nil {
		errhandler.Handler(c, err, "Error updating token")
		return
	}

	// Send response
	res := TaskResponse{
		IsSuccess: true,
		Message:   "success",
	}
	c.JSON(http.StatusOK, res)
}
