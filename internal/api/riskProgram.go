package api

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type riskProgramResponse struct {
	IsSuccess bool `json:"isSuccess"`
	Data riskProgramParams `json:"Data"`
}

type riskProgramParams struct {
	ComputerId int `json:"computerId"`
	RiskLevel int `json:"riskLevel"`
	ProgramId string `json:"programId"`
	ProgramName string `json:"programName"`
	ProgramDate string `json:"programDate"`
	ProgramPath string `json:"programPath"`
}

func riskProgram(c *gin.Context) {

	res := riskProgramResponse{
		IsSuccess: true,
		Data: riskProgramParams{
			ComputerId: 12,
			RiskLevel: 3,
			ProgramId: "lorem",
			ProgramName: "lorem",
			ProgramDate: "05/09/2023",
			ProgramPath: "path/example",
		},
	}

	c.JSON(http.StatusOK, res)
}