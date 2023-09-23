package setting

import (
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/mariadb/query"

	"github.com/gin-gonic/gin"
)

func GetLogs(c *gin.Context) {
	logs, err := query.LoadLogs()
	if err != nil {
		errhandler.Error(c, err, "Error retrieving logs")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"data": logs,
	})
}
