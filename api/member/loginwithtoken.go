package member

import (
	"edetector_API/internal/errhandler"
	"edetector_API/internal/token"
	"edetector_API/pkg/mariadb"
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
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	fmt.Println("Request content: ", req)

	// Verify token
	var verified bool
	userId, err := token.Verify(req.Token)
	if err != nil {
		errhandler.Handler(c, err, "Error verifying token")
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
			errhandler.Handler(c, err, "Error retrieving username")
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
