package userTask

import (
	"edetector_API/pkg/logger"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type taskPacket struct {
	Key      string      `json:"key"`
	Work     string      `json:"work"`
	User     string      `json:"user"`
	Message  interface{} `json:"message"`
}

type taskResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

func ChangeDetectMode(c *gin.Context) {

	req := taskPacket{
		Key:  "8beba472f3f44cabbbb44fd232171933",
		Work: "ChangeDetectMode",
		User: "1",
		Message: "0|0",
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		logger.Error("Error coverting req to JSON: " + err.Error())
	}

	// Set up TCP connection
	conn, err := net.Dial("tcp", "192.168.200.163:1990")
	if err != nil {
		logger.Error("Error connecting to TCP server: " + err.Error())
	}
	defer conn.Close()

	// Send the JSON request
	_, err = conn.Write(reqJSON)
	if err != nil {
		logger.Error("Error sending request: " + err.Error())
	}

	// Receive the response
	responseJSON := make([]byte, 1024)
	_, err = conn.Read(responseJSON)
	if err != nil {
		logger.Error("Error receiving response: " + err.Error())
	}

	// Parse the JSON response
	var response taskResponse
	err = json.Unmarshal(responseJSON, &response)
	if err != nil {
		logger.Error("Error parsing JSON response: " + err.Error())
	}

	fmt.Println("Response:", response.Message)

	res := struct{
		IsSuccess bool `json:"isSuccess"`
		Message string `json:"message"`
	}{
		IsSuccess: response.IsSuccess,
		Message: response.Message,
	}

	c.JSON(http.StatusOK, res)
}