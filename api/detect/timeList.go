package detect

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type timeListParams struct {
	MemoryTimeList []string `json:"memoryTimeList"`
	ForensicsTimeList []string `json:"forensicsTimeList"`
	FileSummaryDownloadTimeList []string `json:"fileSummaryDownloadTimeList"`
}

func TimeList(c *gin.Context) {

	res := Response{
		IsSuccess: true,
		Data: timeListParams{
			MemoryTimeList: []string{"00:00", "03:00", "06:00", "12:00"},
			ForensicsTimeList: []string{"00:00", "03:00", "06:00", "12:00"},
			FileSummaryDownloadTimeList: []string{"00:00", "03:00", "06:00", "12:00"},
		},
	}

	c.JSON(http.StatusOK, res)
}