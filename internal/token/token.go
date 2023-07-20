package token

import (
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ExpirationPeriod = time.Hour * 24

func Generate() (string) {
	return uuid.NewString()
}

func Verify(token string) (int, error) {
	var userId int
	var timestamp sql.NullString
	// check if valid
	query := "SELECT id, token_time FROM user_info WHERE token = ?"
	err := mariadb.DB.QueryRow(query, token).Scan(&userId, &timestamp)
	if err != nil {
		// token incorrect
		if err == sql.ErrNoRows {
			fmt.Println("token " + token + " incorrect")
			logger.Info("token " + token + " incorrect")
			return -1, nil
		} else {
			return -1, err
		}
	}
	if timestamp.Valid {
		// check if token expired
		layout := "2006-01-02 15:04:05"
		parsedTimestamp, err := time.Parse(layout, timestamp.String)
		if err != nil {
			return -1, fmt.Errorf("failed to parse token timestamp: %v", err)
		}
		expirationTime := parsedTimestamp.Add(ExpirationPeriod)
		if time.Now().After(expirationTime) {
			fmt.Println("token " + token + " expired")
			logger.Info("token " + token + " expired")
			return -1, nil
		}
	}
	return userId, nil
}

func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the token from the header
		token := c.GetHeader("Authorization")
		userId, err := Verify(token)
		if err != nil {
			logger.Error("Error verifying token: " + err.Error())
		} else if userId == -1 { // token incorrect
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Set("userID", userId)
			c.Next()
		}
	}
}