package api

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type connectCountParams struct {
	ConnectCount int `json:"ccConnectCount"`
}

func connectCount(c *gin.Context) {

	res := Response{
		IsSuccess: true,
		Data: connectCountParams{
			ConnectCount: 10,
		},
	}

	c.JSON(http.StatusOK, res)
}