package task

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

type TaskRequest struct {
	Key      string      `json:"key"`
	User     string      `json:"user"`
	Message  interface{} `json:"message"`
}

type TaskPacket struct {
	Key      string      `json:"key"`
	Work     string      `json:"work"`
	User     string      `json:"user"`
	Message  interface{} `json:"message"`
}

type TaskResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

// SendTCPRequest sends a TCP request and receives the response
func SendTCPRequest(c *gin.Context, work string, req TaskRequest) (*TaskResponse, error) {
	// Get the TCP connection from the Gin context
	conn, ok := c.Get("tcpConnection")
	if !ok {
		// TCP connection not found in context
		return nil, fmt.Errorf("Failed to get TCP connection")
	}

	// Type assertion to convert the connection to the appropriate type
	tcpConn, ok := conn.(net.Conn)
	if !ok {
		// Invalid TCP connection type
		return nil, fmt.Errorf("Invalid TCP connection")
	}

	// Generate server request
	ServerReq := TaskPacket{
		Key:     req.Key,
		Work:    work,
		User:    req.User,
		Message: req.Message,
	}

	// Convert to JSON
	reqJSON, err := json.Marshal(ServerReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert req to JSON: %s", err.Error())
	}

	// Send the JSON request
	_, err = tcpConn.Write(reqJSON)
	if err != nil {
		return nil, fmt.Errorf("Error sending request to working server: %s", err.Error())
	}

	// Receive the response
	buf := make([]byte, 1024)
	reslen, err := tcpConn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("Error receiving response: %s", err.Error())
	}

	// Parse the JSON response
	var response TaskResponse
	err = json.Unmarshal(buf[:reslen], &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response: %s", err.Error())
	}

	return &response, nil
}