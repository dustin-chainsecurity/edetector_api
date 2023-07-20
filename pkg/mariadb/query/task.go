package query

import (
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"strconv"
)

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
			if count == 0 {  //* no tasks yet
				return -1, 0, nil
			} else { //* all tasks finished
				return 3, 0, nil
			}
		} else {
			return -1, 0, err
		}
	}
	return status, progress, nil //* status: 0, 1, 2
}

func Load_stored_task(taskId string, clientId string, status int) [][]string {
	q := "select task_id, client_id, status from task where "
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

func Update_task_status(taskId string, status int) {
	_, err := mariadb.DB.Exec("update task set status = ? where task_id = ?", status, taskId)
	if err != nil {
		logger.Error("Error updating task status: " + err.Error())
	}
}