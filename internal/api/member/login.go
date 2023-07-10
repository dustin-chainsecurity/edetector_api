package member

import (
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

type UserInfo struct {
	ID       int
	Username string
	Password string
	Token    string
}

func Login(c *gin.Context) {

	// Receive request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request format": err.Error()})
		logger.Error("Invalid request format: " + err.Error())
		return
	}

	var message = "login success"
	var verified = true

	// Get user data
	user_info := UserInfo{
		ID:       -1,
		Username: "Nil",
		Password: "Nil",
		Token:    "Nil",
	}

	query := "SELECT id FROM user WHERE username = ? AND password = MD5(?)"
	err := mariadb.DB.QueryRow(query, req.Username, req.Password).Scan(&user_info.ID)
	if err != nil {
		// Username not exist or wrong password
		if err == sql.ErrNoRows {
			message = "wrong username or password"
			verified = false
			user_info.ID = -1
		} else {
			logger.Error("Error checking user info: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error checking user info": err.Error()})
			return
		}
	}

	if verified {

		user_info.Username = req.Username
		
		// Generate token
		user_info.Token = token.Generate()

		// Update user token
		query = "UPDATE user_info SET token = ?, token_time = CURRENT_TIMESTAMP WHERE id = ?"
		_, err = mariadb.DB.Exec(query, user_info.Token, user_info.ID)
		if err != nil {
			logger.Error("Error updating token: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error updating token": err.Error()})
			return
		}
	}

	// Create response
	res := LoginResponse{
		Success: verified,
		Message: message,
		User: User{
			Username: user_info.Username,
			Token:    user_info.Token,
		},
	}

	c.JSON(http.StatusOK, res)
}