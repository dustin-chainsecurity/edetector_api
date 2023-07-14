package task

import (
	"edetector_API/internal/channel"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/mariadb/query"
	"edetector_API/pkg/redis"
	"encoding/json"
	"time"

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
	query := "INSERT INTO task (task_id, client_id, type, status, progress, timestamp) VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)"
	_, err := mariadb.DB.Exec(query, taskId, deviceId, work, 0, 0)
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

func MonitorStatus() {
	for {
		started_tasks := query.GetStatus("start")
		finished_tasks := query.GetStatus("finish")
		if len(started_tasks) > 0 || len(finished_tasks) > 0 {
			channel.SignalChannel <- "refresh"
			channel.TaskChangeChannel <- started_tasks
			channel.TaskChangeChannel <- finished_tasks
		}
		time.Sleep(5 * time.Second)
	}
}