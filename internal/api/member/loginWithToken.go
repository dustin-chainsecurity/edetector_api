package member

import (
	"edetector_API/pkg/logger"
	"edetector_API/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginWithToken(c *gin.Context) {

	// Receive request
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request format": err.Error()})
		logger.Error("Invalid request format: " + err.Error())
		return
	}

	// Verify token
	verified, err := token.Verify(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error verifying token": err.Error()})
		logger.Error("Error verifying token: " + err.Error())
		return
	}
	var message string
	if verified {
		message = "token verified"
	} else {
		message = "failed"
	}

	// Create response
	var Username = "tmp"
	var Token = "tmp"

	res := LoginResponse{
		Success: verified,
		Message: message,
		Code:    200,
		User: User{
			Username: Username,
			Token:    Token,
		},
	}

	c.JSON(http.StatusOK, res)
}