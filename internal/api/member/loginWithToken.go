package member

import (
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
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
	var verified bool
	userId, err := token.Verify(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error verifying token": err.Error()})
		logger.Error("Error verifying token: " + err.Error())
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
		query := "SELECT username FROM user WHERE id = ?"
		err := mariadb.DB.QueryRow(query, userId).Scan(&username)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error retrieving username": err.Error()})
			logger.Error("Error retrieving username: " + err.Error())
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