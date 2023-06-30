package api

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type timeListResponse struct {
	IsSuccess bool `json:"isSuccess"`
	Data timeListParams `json:"Data"`
}

type timeListParams struct {
	MemoryTimeList []string `json:"memoryTimeList"`
	ForensicsTimeList []string `json:"forensicsTimeList"`
	FileSummaryDownloadTimeList []string `json:"fileSummaryDownloadTimeList"`
}

func timeList(c *gin.Context) {

	res := timeListResponse{
		IsSuccess: true,
		Data: timeListParams{
			MemoryTimeList: []string{"00:00", "03:00", "06:00", "12:00"},
			ForensicsTimeList: []string{"00:00", "03:00", "06:00", "12:00"},
			FileSummaryDownloadTimeList: []string{"00:00", "03:00", "06:00", "12:00"},
		},
	}

	c.JSON(http.StatusOK, res)
}