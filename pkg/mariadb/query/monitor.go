package query

import (
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/logger"
)

func GetStatus(action string) [][]string {
	query := "SELECT client_id, type FROM task WHERE " + action + " = 1"
	var result [][]string
	res, err := mariadb.DB.Query(query)
	if err != nil {
		logger.Error("Error getting status: " + err.Error())
		return result
	}
	defer res.Close()

	for res.Next() {
		tmp := make([]string, 2)
		err := res.Scan(&tmp[0], &tmp[1])
		if err != nil {
			logger.Error("Error scanning task data: " + err.Error())
			return result
		}
		result = append(result, tmp)
	}
	return result
}