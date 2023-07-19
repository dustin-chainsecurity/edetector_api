package saveagent

import (
	"bytes"
	"edetector_API/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Request struct {
	IP         string `json:"ip"`
	Port       string `json:"port"`
	DetectPort string `json:"detect_port"`
}

func SaveAgent(c *gin.Context) {

	request := Request{
		IP:         "192.168.200.161",
		Port:       "5000",
		DetectPort: "5001",
	}

	// Marshal payload into JSON
	payload, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	// Create an HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://192.168.200.167:8080/agent", bytes.NewBuffer(payload)) // need port change
	if err != nil {
		logger.Error("Error creating HTTP request: " + err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending HTTP request: " + err.Error())
		return
	}
	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status code:", response.StatusCode)
		return
	}

	// Create the output file
	agentFile, err := os.Create("agent.exe")
	if err != nil {
		logger.Error("Error creating output file: " + err.Error())
		return
	}
	defer agentFile.Close()

	// Save the response body to the output file
	_, err = io.Copy(agentFile, response.Body)
	if err != nil {
		logger.Error("Error saving response: " + err.Error())
		return
	}

	fmt.Println("Agent saved successfully!")
	res := struct {
		IsSuccess bool `json:"isSuccess"`
		Message string `json:"message"`
	} {
		IsSuccess: true,
		Message: "success",
	}
	c.JSON(http.StatusOK, res)
}