package tcp

import (
	"edetector_API/config"
	Error "edetector_API/internal/error"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

var Conn net.Conn
var err  error

func EstablishTCPConnection() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Establish TCP connection if not already established
		if Conn == nil {
			serverIP := config.Viper.GetString("SERVER_IP")
			serverPort := config.Viper.GetString("SERVER_PORT")
			connectionString := fmt.Sprintf("%s:%s", serverIP, serverPort)
			Conn, err = net.Dial("tcp", connectionString)
			if err != nil {
				Error.Handler(c, err, "Error establishing TCP connection")
				c.Abort()
				return
			}
		}

		// Set the TCP connection in the Gin context for the downstream handlers
		c.Set("tcpConnection", Conn)
		// Proceed to the next handler
		c.Next()
	}
}