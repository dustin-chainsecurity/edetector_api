package member

import (
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
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
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Check if username already exist
	var message string
	var exist bool
	verified := false

	query := "SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)"
	err := mariadb.DB.QueryRow(query, req.Username).Scan(&exist)
	if err != nil {
		logger.Error("Error checking username existence: " + err.Error())
	}

	if !exist {

		// Store user password to db
		query = "INSERT INTO user (username, password) VALUES (?, MD5(?))"
		_, err = mariadb.DB.Exec(query, req.Username, req.Password)
		if err != nil {
			logger.Error("Error storing user data: " + err.Error())
		}

		// Retrieve user id
		var userId int
		query := "SELECT id FROM user WHERE username = ?"
		err := mariadb.DB.QueryRow(query, req.Username).Scan(&userId)
		if err != nil {
			logger.Error("Error retrieving user id: " + err.Error())
		}

		// Storing other user data
		query = "INSERT INTO user_info (id, token, email) VALUES (?, ?, ?)"
		_, err = mariadb.DB.Exec(query, userId, req.Token, req.Email)
		if err != nil {
			logger.Error("Error storing user data: " + err.Error())
		}
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