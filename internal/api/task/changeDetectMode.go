package task

import (
	"edetector_API/pkg/logger"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChangeDetectMode(c *gin.Context, conn net.Conn) {

	req := TaskPacket{
		Key:  "8beba472f3f44cabbbb44fd232171933",
		Work: "ChangeDetectMode",
		User: "1",
		Message: "0|1",
	}

	reqJSON, err := json.Marshal(req)
	if err != nil {
		logger.Error("Error coverting req to JSON: " + err.Error())
	}

	// Send the JSON request
	_, err = conn.Write(reqJSON)
	if err != nil {
		logger.Error("Error sending request: " + err.Error())
	}

	// Change client setting

	// Receive the response
	responseJSON := make([]byte, 1024)
	_, err = conn.Read(responseJSON)
	if err != nil {
		logger.Error("Error receiving response: " + err.Error())
	}

	// Parse the JSON response
	var response TaskResponse
	fmt.Println(response)
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