package member

import (
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type SignUpResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Signup(c *gin.Context) {

	// Receive request
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	// Check if username already exist
	var message string
	verified := false

	exist, err := query.CheckUsername(req.Username)
	if err != nil {
		errhandler.Error(c, err, "Error checking username existence")
		return
	}
	if !exist {
		query.AddUser(req.Username, req.Password, req.Token, req.Email)
		verified = true
		message = "signup success"
	} else {
		verified = false
		message = "username already exist"
	}

	// Create response
	res := SignUpResponse{
		Success: verified,
		Message: message,
	}

	c.JSON(http.StatusOK, res)
}
