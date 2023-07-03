package dashboard

import (
	"net/http"
	"github.com/gin-gonic/gin"	
)

type riskProgramParams struct {
	ComputerId int `json:"computerId"`
	RiskLevel int `json:"riskLevel"`
	ProgramId string `json:"programId"`
	ProgramName string `json:"programName"`
	ProgramDate string `json:"programDate"`
	ProgramPath string `json:"programPath"`
}

func RiskProgram(c *gin.Context) {

	res := Response{
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