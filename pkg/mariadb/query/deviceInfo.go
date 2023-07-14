package query

import (
	"database/sql"
	"edetector_API/pkg/mariadb"
)

type RawDevice struct {
	DeviceID           string
	InnerIP            string
	Network            int
	Process            int
	DeviceName         string
	ScanSchedule       sql.NullString
	ScanFinishTime     sql.NullString
	CollectSchedule    sql.NullString
	CollectFinishTime  sql.NullString
	FileSchedule       sql.NullString
	FileFinishTime     sql.NullString
	ImageFinishTime    sql.NullString
}

func LoadDeviceInfo(deviceId string) (RawDevice, error) {
	var d RawDevice
	query := `
	SELECT C.client_id, C.ip, S.networkreport, S.processreport, I.computername, 
		T.scan_schedule, T.scan_finish_time, T.collect_schedule, T.collect_finish_time, 
		T.file_schedule, T.file_finish_time, T.image_finish_time
	FROM client AS C
	JOIN client_setting AS S ON C.client_id = S.client_id
	JOIN client_info AS I ON S.client_id = I.client_id
	JOIN client_task_status AS T ON I.client_id = T.client_id
	WHERE C.client_id = ?
	`
	err := mariadb.DB.QueryRow(query, deviceId).Scan(
		&d.DeviceID,
		&d.InnerIP,
		&d.Network,
		&d.Process,
		&d.DeviceName,
		&d.ScanSchedule,
		&d.ScanFinishTime,
		&d.CollectSchedule,
		&d.CollectFinishTime,
		&d.FileSchedule,
		&d.FileFinishTime,
		&d.ImageFinishTime,
	)
	if err != nil {
		return d, err
	}

	return d, nil
}

func LoadAllDeviceInfo() ([]RawDevice, error) {
	var devices []RawDevice
	query := `
	SELECT C.client_id, C.ip, S.networkreport, S.processreport, I.computername, 
		T.scan_schedule, T.scan_finish_time, T.collect_schedule, T.collect_finish_time, 
		T.file_schedule, T.file_finish_time, T.image_finish_time
	FROM client AS C
	JOIN client_setting AS S ON C.client_id = S.client_id
	JOIN client_info AS I ON S.client_id = I.client_id
	JOIN client_task_status AS T ON I.client_id = T.client_id
	`
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		return devices, err
	}
	defer rows.Close()

	for rows.Next() {
		var d RawDevice
		err := rows.Scan(
			&d.DeviceID,
			&d.InnerIP,
			&d.Network,
			&d.Process,
			&d.DeviceName,
			&d.ScanSchedule,
			&d.ScanFinishTime,
			&d.CollectSchedule,
			&d.CollectFinishTime,
			&d.FileSchedule,
			&d.FileFinishTime,
			&d.ImageFinishTime,
		)
		if err != nil {
			return devices, err
        }

		devices = append(devices, d)
	}

	return devices, nil
}

func UpdateSchedule(deviceId string, column string, schedule string) error {
	query := "UPDATE client_task_status SET " + column + " = ? WHERE client_id = ?"
	_, err := mariadb.DB.Exec(query, schedule, deviceId)
	if err != nil {
		return err
	}
	return nil
}