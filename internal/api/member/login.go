package member

import (
	"crypto/md5"
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	User    User   `json:"user"`
}

type UserInfo struct {
	ID int
	Username string
	Password string
	Email    string
	Token    string
}

// Expiration time for the token (1 hour)
var tokenExpiration = time.Hour

// Custom Claims structure for JWT
type CustomClaims struct {
	jwt.StandardClaims
}

func Login(c *gin.Context) {

	// Receive request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	var message string
	exist := true
	verified := false

	// Get user data
	user_info := UserInfo{
		ID: -1,
		Username: "Nil",
		Password: "Nil",
		Email: "Nil",
		Token: "Nil",
	}

	query := "SELECT id, password FROM user WHERE username = ?"
	err := mariadb.DB.QueryRow(query, req.Username).Scan(&user_info.ID, &user_info.Password)
	if err != nil {
		// Username not exist
		if err == sql.ErrNoRows {
			exist = false
			message = "username not exist"
		}
		logger.Error("Error retrieving password: " + err.Error())
	}

	if exist {
		// Check user password
		hash := md5.Sum([]byte(req.Password))
		encoded := hex.EncodeToString(hash[:])
		if encoded == user_info.Password {
			message = "login success"
			verified = true
		} else {
			message = "password incorrect"
			verified = false
			user_info.ID = -1
			user_info.Username = "Nil"
		}

		if verified {
			// Create a new token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(tokenExpiration).Unix(),
				},
			})

			// Sign the token with your secret key
			var jwtSecret = []byte(user_info.Username)
			tokenString, err := token.SignedString(jwtSecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}

			// Set the token as a cookie
			cookie := &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				Expires:  time.Now().Add(tokenExpiration),
				HttpOnly: true,
			}
			http.SetCookie(c.Writer, cookie)
			user_info.Token = tokenString

			// Update user token
			query = "UPDATE user_info SET token = ? WHERE id = ?"
			_, err = mariadb.DB.Exec(query, tokenString, user_info.ID)
			if err != nil {
				logger.Error("Error updating token: " + err.Error())
			}

			// Get user_info data
			query = "SELECT email FROM user_info WHERE id = ?"
			err = mariadb.DB.QueryRow(query, user_info.ID).Scan(&user_info.Email)
			if err != nil {
				logger.Error("Error retrieving user_info: " + err.Error())
			}
		}
	}

	// Create response
	res := LoginResponse{
		Success: verified,
		Message: message,
		Code:    200,
		User: User{
			ID:       user_info.ID,
			Username: user_info.Username,
			Email:    user_info.Email,
			Token:    user_info.Token,
		},
	}

	c.JSON(http.StatusOK, res)
}

// INSERT INTO user (username, password) VALUES ('example_username', 'example_password');
// INSERT INTO user_info (id, token, email) VALUES (1, 'loremipsumdolorsitamet', 'chiehyu@exampe.com');
