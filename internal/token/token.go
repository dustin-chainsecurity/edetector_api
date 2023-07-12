package token

import (
	"database/sql"
	"edetector_API/pkg/mariadb"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ExpirationPeriod = time.Hour * 24

func Generate() (string) {
	return uuid.NewString()
}

func Verify(token string) (int, error) {

	var userId = -1
	var timestamp string

	// check if valid
	query := "SELECT id, token_time FROM user_info WHERE token = ?"
	err := mariadb.DB.QueryRow(query, token).Scan(&userId, &timestamp)
	if err != nil {
		// token not exist
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("token not exist")
		} else {
			return -1, fmt.Errorf(err.Error())
		}
	}

	// check if token expired
	layout := "2006-01-02 15:04:05"
	parsedTimestamp, err := time.Parse(layout, timestamp)
	if err != nil {
		return -1, fmt.Errorf("failed to parse token timestamp: %v", err)
	}
	expirationTime := parsedTimestamp.Add(ExpirationPeriod)
	if time.Now().After(expirationTime) {
		return -1, fmt.Errorf("token expired")
	}

	return userId, nil
}