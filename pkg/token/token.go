package token

import (
	"crypto/rand"
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

func Verify(token string) (bool, error) {

	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid token format")
	}
	tokenTime, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("failed to parse token time: %v", err)
	}
	expiration := time.Unix(tokenTime, 0).Add(ExpirationPeriod)
	if time.Now().After(expiration) {
		return false, fmt.Errorf("token expired")
	}
	return true, nil
}

func generateRandomTokenValue() (string, error) {

	tokenValueBytes := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(tokenValueBytes)
	if err != nil {
		return "", err
	}
	// convert to base64
	tokenValue := base64.URLEncoding.EncodeToString(tokenValueBytes)
	return tokenValue, nil
}