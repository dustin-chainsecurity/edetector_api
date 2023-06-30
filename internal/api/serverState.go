package api

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type serverStateResponse struct {
	IsSuccess bool `json:"isSuccess"`
	Data serverParams `json:"Data"`
}

type serverParams struct {
	CPU int `json:"cpu"`
	RAMALL int `json:"RAMALL"`
	RAMUsed int `json:"RAMUsed"`
	Storage int `json:"storage"`
}

func serverState(c *gin.Context) {

	res := serverStateResponse{
		IsSuccess: true,
		Data: serverParams{
			CPU: 98,
			RAMALL: 1024,
			RAMUsed: 1024,
			Storage: 50,
		},
	}

	c.JSON(http.StatusOK, res)
}