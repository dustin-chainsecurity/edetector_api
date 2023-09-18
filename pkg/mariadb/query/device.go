package query

import (
	"crypto/rand"
	"database/sql"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"encoding/json"
	"fmt"
)

type RawDevice struct {
	DeviceID           string
	InnerIP            string
	Mac                string
	Network            int
	Process            int
	DeviceName         string
	Version            string
	ScanSchedule       sql.NullString
	ScanFinishTime     sql.NullString
	CollectSchedule    sql.NullString
	CollectFinishTime  sql.NullString
	FileSchedule       sql.NullString
	FileFinishTime     sql.NullString
	ImageFinishTime    sql.NullString
}

type ClientOnlineStatus struct {
	Status int
	Time   string
}

func LoadDeviceInfo(deviceId string) (RawDevice, error) {
	var d RawDevice
	query := `
	SELECT C.client_id, C.ip, C.mac, S.networkreport, S.processreport, I.computername, I.fileversion,
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
		&d.Mac,
		&d.Network,
		&d.Process,
		&d.DeviceName,
		&d.Version,
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
	SELECT C.client_id, C.ip, C.mac, S.networkreport, S.processreport, I.computername, I.fileversion, 
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
			&d.Mac,
			&d.Network,
			&d.Process,
			&d.DeviceName,
			&d.Version,
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

func GetDetectMode(deviceId string) (int, int, error) {
	var process, network int
	query := "SELECT processreport, networkreport FROM client_setting WHERE client_id = ?"
	err := mariadb.DB.QueryRow(query, deviceId).Scan(&process, &network)
	if err != nil {
		return 0, 0, err
	}
	return process, network, nil
}

func UpdateDetectMode(mode int, deviceId string) error {
	query := "UPDATE client_setting SET processreport = ?, networkreport = ? WHERE client_id = ?"
	_, err := mariadb.DB.Exec(query, mode, mode, deviceId)
	return err
}

func CheckDevice(deviceId string) (bool, error) {
	query := `
		SELECT EXISTS(SELECT client_id FROM client WHERE client_id = ?);
	`
	var exist bool
	err := mariadb.DB.QueryRow(query, deviceId).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func AddDevice(id string, name string, ip string) error {
	// check if device exist
	if exist, err := CheckDevice(id); err != nil {
		return err
	} else if exist {
		return fmt.Errorf("device " + id + " already exist")
	}
	query := `
		INSERT INTO client (client_id, ip, mac) VALUES (?, ?, ?);
	`
	_, err := mariadb.DB.Exec(query, id, ip, generateMACAddress())
	if err != nil {
		return err
	}
	query = `
		INSERT INTO client_setting (client_id, networkreport, processreport) VALUES (?, 0, 0);
	`
	_, err = mariadb.DB.Exec(query, id)
	if err != nil {
		return err
	}
	query = `
		INSERT INTO client_info (client_id, sysinfo, osinfo, computername, username, fileversion, boottime) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err = mariadb.DB.Exec(query, id, "x64", "Windows 10 Home", name, "SYSTEM", "3.4.2.0,1988,1989", "20230706110846")
	if err != nil {
		return err
	}
	query = `
		INSERT INTO client_task_status (client_id) VALUES (?);
	`
	_, err = mariadb.DB.Exec(query, id)
	if err != nil {
		return err
	}
	onlineStatusInfo := ClientOnlineStatus {
		Status: 0,
		Time:   "20230706110846",
	}
	redis.Redis_set(id, onlineStatusInfo.Marshal())	
	return nil
}

func (c *ClientOnlineStatus) Marshal() string {
	json, err := json.Marshal(c)
	if err != nil {
		logger.Error("Error in json marshal" + err.Error())
	}
	return string(json)
}

func generateMACAddress() string {
	mac := make([]byte, 6)
	_, err := rand.Read(mac)
	if err != nil {
		panic(err)
	}
	mac[0] |= 2
	macAddress := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])

	return macAddress
}