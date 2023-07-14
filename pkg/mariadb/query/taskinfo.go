package query

import (
	"database/sql"
	"edetector_API/pkg/mariadb"
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