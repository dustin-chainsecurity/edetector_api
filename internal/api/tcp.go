package api

import (
	"edetector_API/config"
	"edetector_API/pkg/logger"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EstablishTCPConnection() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Establish TCP connection if not already established
		if conn == nil {
			serverIP := config.Viper.GetString("SERVER_IP")
			serverPort := config.Viper.GetString("SERVER_PORT")
			connectionString := fmt.Sprintf("%s:%s", serverIP, serverPort)
			conn, err = net.Dial("tcp", connectionString)
			if err != nil {
				logger.Error("Error setting up TCP connection: " + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to establish TCP connection"})
				c.Abort()
				return
			}
		}

		// Set the TCP connection in the Gin context for the downstream handlers
		c.Set("tcpConnection", conn)
		// Proceed to the next handler
		c.Next()
	}
}