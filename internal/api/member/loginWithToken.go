package member

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginWithToken(c *gin.Context) {

	// Receive request
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// verify token

}