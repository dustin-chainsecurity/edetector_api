package task

import (
	"edetector_API/pkg/logger"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartCollect(c *gin.Context) {

	// Get the TCP connection from the Gin context
	conn, ok := c.Get("tcpConnection")
	if !ok {
		// TCP connection not found in context
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get TCP connection"})
		return
	}

	// Type assertion to convert the connection to the appropriate type
	tcpConn, ok := conn.(net.Conn)
	if !ok {
		// Invalid TCP connection type
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid TCP connection"})
		return
	}

	// Receive user request
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Generate server request
	ServerReq := TaskPacket{
		Key:  req.Key,
		Work: "StartCollect",
		User: req.User,
		Message: req.Message,
	}

	// Convert to JSON
	reqJSON, err := json.Marshal(ServerReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed to convert req to JSON": err.Error()})
		logger.Error("Error coverting req to JSON: " + err.Error())
		return
	}

	// Send the JSON request
	_, err = tcpConn.Write(reqJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error sending request to working server": err.Error()})
		logger.Error("Error sending request: " + err.Error())
		return
	}

	// Functioning

	// Receive the response
	buf := make([]byte, 1024)
	reslen, err := tcpConn.Read(buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error receiving response": err.Error()})
		logger.Error("Error receiving response: " + err.Error())
		return
	}

	// Parse the JSON response
	var response TaskResponse
	err = json.Unmarshal(buf[:reslen], &response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error parsing JSON response": err.Error()})
		logger.Error("Error parsing JSON response: " + err.Error())
		return
	}
	fmt.Println("Response:", response.Message)

	res := TaskResponse {
		IsSuccess: response.IsSuccess,
		Message: response.Message,
	}

	c.JSON(http.StatusOK, res)	
}