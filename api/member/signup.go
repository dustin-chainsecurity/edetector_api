package member

import (
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
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
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	// Check if username already exist
	var message string
	var exist bool
	verified := false

	query := "SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)"
	err := mariadb.DB.QueryRow(query, req.Username).Scan(&exist)
	if err != nil {
		errhandler.Handler(c, err, "Error checking username existence")
		return
	}

	if !exist {

		// Store user password to db
		query = "INSERT INTO user (username, password) VALUES (?, MD5(?))"
		_, err = mariadb.DB.Exec(query, req.Username, req.Password)
		if err != nil {
			errhandler.Handler(c, err, "Error storing user data")
			return
		}

		// Retrieve user id
		var userId int
		query := "SELECT id FROM user WHERE username = ?"
		err := mariadb.DB.QueryRow(query, req.Username).Scan(&userId)
		if err != nil {
			errhandler.Handler(c, err, "Error retrieving user id")
			return
		}

		// Storing other user data
		query = "INSERT INTO user_info (id, token, email) VALUES (?, ?, ?)"
		_, err = mariadb.DB.Exec(query, userId, req.Token, req.Email)
		if err != nil {
			errhandler.Handler(c, err, "Error storing user data")
			return
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
