package task

import (
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"encoding/json"
	"github.com/google/uuid"
)

type TaskPacket struct {
	Key      string      `json:"key"`
	Work     string      `json:"work"`
	User     string      `json:"user"`
	Message  string      `json:"message"`
}

type TaskResponse struct {
	IsSuccess  bool   `json:"isSuccess"`
	Message    string `json:"message"`
}

func AddTask(deviceId string, work string, msg string) error {

	taskId := uuid.NewString()

	// store into mariaDB
	query := "INSERT INTO task (task_id, client_id, type, status, start, finish, timestamp) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)"
	_, err := mariadb.DB.Exec(query, taskId, deviceId, work, 0, 0, 0)
	if err != nil {
		return err
	}

	// store into redis
	pkt := TaskPacket{
		Key: deviceId,
		Work: work,
		User: "1",
		Message: msg,
	}
	pktString, err := json.Marshal(pkt)
	if err != nil {
		return err
	}

	err = redis.Redis_set(taskId, string(pktString))
	return err
}