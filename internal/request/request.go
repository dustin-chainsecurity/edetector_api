package request

import (
	"bytes"
	"edetector_API/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

type taskRequest struct {
	TaskID  string  `json:"taskId"`
}

type updateRequest struct {
	DeviceID  string  `json:"deviceId"`
}

type sendMissionRequest struct {
	Action    string    `json:"action"`
	Devices   []string  `json:"deviceId"`
}

func SendToServer(id string) {
	request := taskRequest {
		TaskID: id,
	}
	SendRequest(request, "http://192.168.200.163:5000/sendTask")
}

func SendUpdateTaskToApi(id string) {
	request := updateRequest {
		DeviceID: id,
	}
	SendRequest(request, "http://192.168.200.161:5050/updateTask")
}

func SendMissionToApi(work string, devices []string) {
	request := sendMissionRequest {
		Action:  work,
		Devices: devices,
	}
	SendRequest(request, "http://192.168.200.161:5000/sendMission")
}

func SendRequest(request interface{}, path string) {
	// Marshal payload into JSON
	payload, err := json.Marshal(request)
	if err != nil {
		logger.Error("Error marshaling JSON: " + err.Error())
		return
	}
	// Create an HTTP request
	client := &http.Client{}
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(payload))
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
}