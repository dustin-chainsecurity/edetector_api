package api

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type riskComputerParams struct {
	ComputerId int `json:"computerId"`
	RiskLevel int `json:"riskLevel"`
}

func riskComputer(c *gin.Context) {

	res := Response{
		IsSuccess: true,
		Data: riskComputerParams{
			ComputerId: 12,
			RiskLevel: 3,
		},
	}

	c.JSON(http.StatusOK, res)
}