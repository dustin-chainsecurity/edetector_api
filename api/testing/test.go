package testing

import (
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	id := c.Param("id")
	if id == "1" {
		c.JSON(400, gin.H{
			"message": "nil",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}