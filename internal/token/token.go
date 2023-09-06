package token

import (
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ExpirationPeriod = time.Hour * 24

func Generate() string {
	return uuid.NewString()
}

func VerifyAdmin(token string) (int, error) {
	return 0, nil
	//ToDo
}

func Verify(token string) (int, error) {
	userID, err := query.CheckToken(token, ExpirationPeriod)
	return userID, err
}

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the token from the header
		token := c.GetHeader("Authorization")
		if token == "" {
			errhandler.Unauthorized(c, fmt.Errorf("token not provided"), "Error verifying token")
			c.Abort()
			return
		}
		userId, err := Verify(token)
		if err != nil {
			errhandler.Error(c, err, "Error verifying token")
			c.Abort()
			return
		} else if userId == -1 { // token incorrect
			errhandler.Unauthorized(c, fmt.Errorf("token incorrect"), "Error verifying token")
			c.Abort()
			return
		} else {
			c.Set("userID", userId)
			c.Next()
		}
	}
}

func TokenAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the token from the header
		token := c.GetHeader("Authorization")
		userId, err := VerifyAdmin(token)
		if err != nil {
			errhandler.Error(c, err, "Error verifying token")
			c.Abort()
			return
		} else if userId == -1 { // token incorrect
			errhandler.Unauthorized(c, fmt.Errorf("token incorrect"), "Error verifying token")
			c.Abort()
			return
		} else {
			c.Set("userID", userId)
			c.Next()
		}
	}
}
