package login

import (
	"net/http"
	"github.com/gin-gonic/gin"
	_ "crypto/md5"
	_ "encoding/hex"
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

func Handle(c *gin.Context) {

	// Receive request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check password
	

	// Create response
	res := LoginResponse{
		Success: true,
		Message: "登入成功",
		Code:    200,
		User: User{
			ID:       123,
			Username: req.Username,
			Email:    "john_doe@example.com",
			Token:    "anyScriptyouchooseorLoremLorem",
		},
	}

	c.JSON(http.StatusOK, res)
}