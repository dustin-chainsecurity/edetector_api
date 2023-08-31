package query

import (
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"fmt"
	"strconv"
)

func AddTask(taskId string, clientId string, work string) error {
	query := "INSERT INTO task (task_id, client_id, type, status, progress, timestamp) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)"
	_, err := mariadb.DB.Exec(query, taskId, clientId, work, 0, 0)
	if err != nil {
		return err
	}
	return nil
}

func LoadTaskStatus(deviceId string, work string) (int, int, error) {
	var status, progress int
	query := "SELECT status, progress FROM task WHERE client_id = ? AND type = ? AND status != 3"
	err := mariadb.DB.QueryRow(query, deviceId, work).Scan(&status, &progress)
	if err != nil {
		if err == sql.ErrNoRows {
			var count int
			query = "SELECT COUNT(*) FROM task WHERE client_id = ? AND type = ? AND status = 3"
			err = mariadb.DB.QueryRow(query, deviceId, work).Scan(&count)
			if err != nil {
				return -1, 0, err
			}
			if count == 0 {                 //* no task yet
				return -1, 0, nil
			} else {                        //* task finished
				return 3, 0, nil
			}
		} else {
			return -1, 0, err
		}
	}
	return status, progress, nil            //* status: 0, 1, 2
}

func LoadStoredTask(taskId string, clientId string, status int) [][]string {
	q := "SELECT task_id, client_id, status FROM task WHERE "
	var result [][]string
	if clientId != "nil" {
		q = q + "client_id = " + clientId
	}
	if taskId != "nil" {
		q = q + "task_id = " + taskId
	}
	if status != -1 {
		q = q + "status = " + strconv.Itoa(status)
	}
	res, err := mariadb.DB.Query(q)
	if err != nil {
		logger.Error("Error loading tasks: " + err.Error())
		return result
	}
	defer res.Close()
	l, _ := res.Columns()
	for res.Next() {
		tmp := make([]string, len(l))
		err := res.Scan(&tmp[0], &tmp[1], &tmp[2])
		if err != nil {
			logger.Error("Error loading tasks: " + err.Error())
			return result
		}
		result = append(result, tmp)
	}
	return result
}

func CheckProcessingTask(deviceId string, work string) (bool, error) {
	var exist bool
	query := "SELECT EXISTS(SELECT task_id FROM task WHERE client_id = ? AND type = ? AND status != 3 AND status != 4)"
	err := mariadb.DB.QueryRow(query, deviceId, work).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func CheckTask(TaskId string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM task WHERE task_id = ?"
	err := mariadb.DB.QueryRow(query, TaskId).Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, fmt.Errorf("task " + TaskId + " not exist")
	} else {
		return true, nil
	}
}

func UpdateTaskStatus(taskId string, status int) {
	_, err := mariadb.DB.Exec("UPDATE task SET status = ? WHERE task_id = ?", status, taskId)
	if err != nil {
		logger.Error("Error updating task status: " + err.Error())
	}
}

func UpdateTaskProgress(taskId string, progress int) error {
	_, err := mariadb.DB.Exec("UPDATE task SET progress = ? WHERE task_id = ?", progress, taskId)
	if err != nil {
		return err
	}
	return nil
}
