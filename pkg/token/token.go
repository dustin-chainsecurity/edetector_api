package token

import (
	"crypto/rand"
	"database/sql"
	"edetector_API/pkg/mariadb"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ExpirationPeriod = time.Hour * 24

func Generate() (string, error) {
	tokenValue, err := generateRandomTokenValue()
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("%s:%d", tokenValue, time.Now().Unix())
	return token, nil
}

func Verify(token string) (int, error) {

	var userId = -1

	// check token format
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return userId, fmt.Errorf("invalid token format")
	}
	tokenTime, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return userId, fmt.Errorf("failed to parse token time: %v", err)
	}

	// check if expired
	expiration := time.Unix(tokenTime, 0).Add(ExpirationPeriod)
	if time.Now().After(expiration) {
		return userId, fmt.Errorf("token expired")
	}

	// check if valid
	query := "SELECT id FROM user_info WHERE token = ?"
	err = mariadb.DB.QueryRow(query, token).Scan(&userId)
	if err != nil {
		// token not exist
		if err == sql.ErrNoRows {
			return userId, fmt.Errorf("token not exist")
		} else {
			return userId, fmt.Errorf(err.Error())
		}
	}

	return userId, nil
}

func generateRandomTokenValue() (string, error) {
	//change to uuid token
	tokenValueBytes := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(tokenValueBytes)
	if err != nil {
		return "", err
	}
	// convert to base64
	tokenValue := base64.URLEncoding.EncodeToString(tokenValueBytes)
	return tokenValue, nil
}
