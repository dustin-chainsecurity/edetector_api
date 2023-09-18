package query

import "edetector_API/pkg/mariadb"

type Log struct {
	Id        int    `json:"id"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

func LoadLogs() ([]Log, error) {
	query := "SELECT id, timestamp, level, service, content FROM log ORDER BY timestamp DESC"
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []Log
	for rows.Next() {
		var log Log
		if err := rows.Scan(&log.Id, &log.Timestamp, &log.Level, &log.Service, &log.Content); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
