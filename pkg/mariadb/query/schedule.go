package query

import (
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
)

func LoadStoredSchedule() [][]sql.NullString {
	query := "select client_id, scan_schedule, collect_schedule, file_schedule from client_task_status"
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		logger.Error("Error loading schedule: " + err.Error())
		return [][]sql.NullString{}
	}
	defer rows.Close()

	result := [][]sql.NullString{}
	for rows.Next() {
		tmp := make([]sql.NullString, 4)
		err := rows.Scan(&tmp[0], &tmp[1], &tmp[2], &tmp[3])
		if err != nil {
			logger.Error("Error loading schedule: " + err.Error())
			return [][]sql.NullString{}
		}
		result = append(result, tmp)
	}
	return result
}

func UpdateSchedule(deviceId string, column string, schedule string, mode bool) error {
	query := "UPDATE client_task_status SET " + column + " = "
	if mode {
		query += schedule  
	} else {
		query += "NULL"
	}
	query += " WHERE client_id = ?"
	_, err := mariadb.DB.Exec(query, deviceId)
	if err != nil {
		return err
	}
	return nil
}