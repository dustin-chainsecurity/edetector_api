package task

import (
	"bytes"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"encoding/json"
	"fmt"
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
	query := "INSERT INTO task (task_id, client_id, status, timestamp) VALUES (?, ?, ?, CURRENT_TIMESTAMP)"
	_, err := mariadb.DB.Exec(query, taskId, deviceId, 0)
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

	// fill to 1024 bytes
	if len(pktString) >= 1024 {
		return fmt.Errorf("task packet too long")
	} else {
		zeroBytes := bytes.Repeat([]byte{0}, 1024 - len(pktString))
		pktString = append(pktString, zeroBytes...)
	}

	err = redis.Redis_set(taskId, string(pktString))
	return err
}