package task

import (
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
	// check processing tasks
	if exist, err := query.CheckProcessingTask(deviceId, work); err != nil {
		return "", err
	} else if exist {
		return "", fmt.Errorf("already processing the same task: %s for device: %s", work, deviceId)
	}
	// generate taskid
	taskId := uuid.NewString()
	// store into mariaDB
	err := query.AddTask(taskId, deviceId, work)
	if err != nil {
		return "", err
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
		return "", err
	}
	err = redis.Redis_set(taskId, string(pktString))
	return taskId, err
}

func processSchedule(date int, time int) string {
	return fmt.Sprintf("%d|%d", date, time)
}