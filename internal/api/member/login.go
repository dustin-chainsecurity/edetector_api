package member

import (
	"crypto/md5"
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/token"
	"encoding/hex"
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
	Code    int    `json:"code"`
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

	var message string
	exist := true
	verified := false

	// Get user data
	user_info := UserInfo{
		ID:       -1,
		Username: "Nil",
		Password: "Nil",
		Token:    "Nil",
	}

	query := "SELECT password FROM user WHERE username = ?"
	err := mariadb.DB.QueryRow(query, req.Username).Scan(&user_info.Password)
	if err != nil {
		// Username not exist
		if err == sql.ErrNoRows {
			exist = false
			message = "username not exist"
		} else {
			logger.Error("Error retrieving password: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error retrieving password": err.Error()})
			return
		}
	}

	if exist {
		// Check user password
		hash := md5.Sum([]byte(req.Password))
		encoded := hex.EncodeToString(hash[:])
		if encoded == user_info.Password {
			message = "login success"
			verified = true
		} else {
			message = "password incorrect"
			verified = false
			user_info.ID = -1
			user_info.Username = "Nil"
		}

		if verified {
			// Generate token
			token, err := token.Generate()
			if err != nil {
				logger.Error("Error generating token: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"Error generating token": err.Error()})
				return
			}
			user_info.Token = token

			// Update user token
			query = "UPDATE user_info SET token = ? WHERE id = ?"
			_, err = mariadb.DB.Exec(query, token, user_info.ID)
			if err != nil {
				logger.Error("Error updating token: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"Error updating token": err.Error()})
				return
			}
		}
	}

	// Create response
	res := LoginResponse{
		Success: verified,
		Message: message,
		Code:    200,
		User: User{
			Username: user_info.Username,
			Token:    user_info.Token,
		},
	}

	c.JSON(http.StatusOK, res)
}

// INSERT INTO user (username, password) VALUES ('example_username', 'example_password');
// INSERT INTO user_info (id, token, email) VALUES (1, 'loremipsumdolorsitamet', 'chiehyu@exampe.com');
