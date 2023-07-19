package task

import (
	"edetector_API/pkg/mariadb"
	"edetector_API/internal/Error"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateProgress(c *gin.Context) {

	// Receive request
	var req struct {
		TaskId    string  `json:"taskId"`
		Progress  int     `json:"progress"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Error.Handler(c, err, "Invalid request format")
		return
	}

	query := "UPDATE task SET progress = ? WHERE task_id = ?"
	_, err := mariadb.DB.Exec(query, req.Progress, req.TaskId)
	if err != nil {
		Error.Handler(c, err, "Error updating token")
		return
	}

	// Send response
	res := TaskResponse {
		IsSuccess: true,
		Message: "success",
	}
	c.JSON(http.StatusOK, res)
}