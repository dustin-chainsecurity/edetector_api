package task

import (
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/mariadb/query"
	"edetector_API/pkg/redis"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type TaskPacket struct {
	Key     string `json:"key"`
	Work    string `json:"work"`
	User    string `json:"user"`
	Message string `json:"message"`
}

type TaskResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

type ScheduleScanRequest struct {
	Mode    bool     `json:"mode"`
	Time    string   `json:"time"`
	Devices []string `json:"deviceId"`
}

type ScheduleRequest struct {
	Mode    bool     `json:"mode"`
	Date    int      `json:"date"`
	Time    int      `json:"time"`
	Devices []string `json:"deviceId"`
}

func addTask(deviceId string, work string, msg string) (string, error) {

	// check deviceId
	exist, err := query.CheckDevice(deviceId)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", fmt.Errorf("deviceId not exist")
	}
	// generate taskid
	taskId := uuid.NewString()
	// store into mariaDB
	query := "INSERT INTO task (task_id, client_id, type, status, progress, timestamp) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)"
	_, err = mariadb.DB.Exec(query, taskId, deviceId, work, 0, 0)
	if err != nil {
		return taskId, err
	}
	// store into redis
	pkt := TaskPacket{
		Key:     deviceId,
		Work:    work,
		User:    "1",
		Message: msg,
	}
	pktString, err := json.Marshal(pkt)
	if err != nil {
		return taskId, err
	}
	err = redis.Redis_set(taskId, string(pktString))
	return taskId, err
}

func processSchedule(date int, time int) string {
	return fmt.Sprintf("%d|%d", date, time)
}