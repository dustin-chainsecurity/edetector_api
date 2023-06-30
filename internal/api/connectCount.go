package api

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type connectCountResponse struct {
	IsSuccess bool `json:"isSuccess"`
	Data connectCountParams `json:"Data"`
}

type connectCountParams struct {
	ConnectCount int `json:"ccConnectCount"`
}

func connectCount(c *gin.Context) {

	res := connectCountResponse{
		IsSuccess: true,
		Data: connectCountParams{
			ConnectCount: 10,
		},
	}

	c.JSON(http.StatusOK, res)
}