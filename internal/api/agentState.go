package api

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type agentStateResponse struct {
	IsSuccess bool `json:"isSuccess"`
	Data agentParams `json:"Data"`
}

type agentParams struct {
	AgentID int `json:"agentId"`
	AgentIsConnect bool `json:"agentIsConnect"`
}

func agentState(c *gin.Context) {

	res := agentStateResponse{
		IsSuccess: true,
		Data: agentParams{
			AgentID: 12,
			AgentIsConnect: true,
		},
	}

	c.JSON(http.StatusOK, res)
}