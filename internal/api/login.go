package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"crypto/md5"
	"encoding/hex"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/logger"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	User    User   `json:"user"`
}

type UserInfo struct {
	ID int
	Username string
	Password string
	Email    string
	Token    string
}

func login(c *gin.Context) {

	// Receive request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user data
	stmt, err := mariadb.DB.Prepare("SELECT id, username, password FROM user WHERE username = ?")
	if err != nil {
		logger.Error("Error preparing stmt to get user data: " + err.Error())
	}
	row := stmt.QueryRow(req.Username)
	user_data := UserInfo{}
	err = row.Scan(&user_data.ID, &user_data.Username, &user_data.Password)
	if err != nil {
		logger.Error("Error scanning user data: " + err.Error())
	}

	// Get user_info data
	stmt, err = mariadb.DB.Prepare("SELECT token, email FROM user_info WHERE id = ?")
	if err != nil {
		logger.Error("Error preparing stmt to get user_info data: " + err.Error())
	}
	row = stmt.QueryRow(user_data.ID)
	err = row.Scan(&user_data.Token, &user_data.Email)
	if err != nil {
		logger.Error("Error scanning user_info data: " + err.Error())
	}

	// Check user password
	var message string
	var verified bool
	hash := md5.Sum([]byte(req.Password))
	encoded := hex.EncodeToString(hash[:])
	if encoded == user_data.Password {
		message = "login successed"
		verified = true
	} else {
		message = "login error"
		verified = false
	}

	// Create response
	res := LoginResponse{
		Success: verified,
		Message: message,
		Code:    200,
		User: User{
			ID:       user_data.ID,
			Username: req.Username,
			Email:    user_data.Email,
			Token:    user_data.Token,
		},
	}

	c.JSON(http.StatusOK, res)
}

// INSERT INTO user (username, password) VALUES ('example_username', 'example_password');
// INSERT INTO user_info (id, token, email) VALUES (1, 'loremipsumdolorsitamet', 'chiehyu@exampe.com');