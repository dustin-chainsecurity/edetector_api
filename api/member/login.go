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

var err error

func Login(c *gin.Context) {

	// Receive request
	var req LoginRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	var message = "login success"
	var verified = true

	// Get user data
	user_info := UserInfo{
		ID:       -1,
		Username: "Nil",
		Password: "Nil",
		Token:    "Nil",
	}

	user_info.ID, err = query.CheckUser(req.Username, req.Password)
	if err != nil {
		errhandler.Handler(c, err, "Error checking user info")
		return
	}
	if user_info.ID == -1 {
		message = "wrong username or password"
		verified = false
	}

	if verified {
		user_info.Username = req.Username
		// Generate token
		user_info.Token = token.Generate()
		// Update user token
		err = query.UpdateToken(user_info.Token, user_info.ID)
		if err != nil {
			errhandler.Handler(c, err, "Error updating token")
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
