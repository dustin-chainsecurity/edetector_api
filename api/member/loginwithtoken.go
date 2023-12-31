package member

import (
	"edetector_API/internal/token"
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginWithToken(c *gin.Context) {

	// Receive request
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	// Verify token
	var verified bool
	userId, err := token.Verify(req.Token)
	if err != nil {
		errhandler.Error(c, err, "Error verifying token")
		return
	}
	var message, username, token string
	if userId == -1 {
		verified = false
		message = "failed"
		username = "Nil"
		token = "Nil"
	} else {
		verified = true
		message = "token verified"
		token = req.Token
		// Get username
		username, err = query.GetUsername(userId)
		if err != nil {
			errhandler.Error(c, err, "Error retrieving username")
			return
		}
	}

	res := LoginResponse{
		Success: verified,
		Message: message,
		User: User{
			Username: username,
			Token:    token,
		},
	}

	c.JSON(http.StatusOK, res)
}
