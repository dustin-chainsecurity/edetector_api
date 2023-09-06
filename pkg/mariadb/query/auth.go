package query

import (
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"time"
)

func CheckToken(token string, expirationPeriod time.Duration) (int, error){
	var userId int
	var timestamp sql.NullString
	query := "SELECT id, token_time FROM user_info WHERE token = ?"
	err := mariadb.DB.QueryRow(query, token).Scan(&userId, &timestamp)
	if err != nil {
		// token incorrect
		if err == sql.ErrNoRows {
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
			return -1, err
		}
		expirationTime := parsedTimestamp.Add(expirationPeriod)
		if time.Now().After(expirationTime) {
			logger.Info("token " + token + " expired")
			return -1, nil
		}
	}
	return userId, nil
}