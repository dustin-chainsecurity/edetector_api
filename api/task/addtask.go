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
	User    int    `json:"user"`
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

func addTask(userId int, deviceId string, work string, msg string, imagetype string) (string, error) {
	// check processing tasks
	if exist, err := query.CheckProcessingTask(deviceId, work); err != nil {
		return "", err
	} else if exist {
		return "", fmt.Errorf("already processing the same task: %s for device: %s", work, deviceId)
	}
	// check terminate status
	if exist, err := query.CheckTerminateStatus(deviceId); err != nil {
		return "", err
	} else if exist {
		return "", fmt.Errorf("device: %s is still terminating", deviceId)
	}
	// generate taskid
	taskId := uuid.NewString()
	// store into mariaDB
	err := query.AddTask(taskId, deviceId, work)
	if err != nil {
		return "", err
	}
	// store into redis
	var pkt TaskPacket
	if work == "StartGetImage" {
		pkt = TaskPacket{
			Key:     deviceId,
			Work:    work,
			User:	 userId,
			Message: imagetype,
		}
	} else {
		pkt = TaskPacket{
			Key:     deviceId,
			Work:    work,
			User:    userId,
			Message: msg,
		}
	}
	pktString, err := json.Marshal(pkt)
	if err != nil {
		return "", err
	}
	err = redis.Redis_set(taskId, string(pktString))
	return taskId, err
}