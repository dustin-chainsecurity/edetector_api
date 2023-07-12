package packet

import (
	"database/sql"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"	
)

type dateForm struct {
	Date   string    `json:"date"`
	Time   string    `json:"time"`
}

type processing struct {
	IsFinish    bool   `json:"isFinish"`
	Progress    int    `json:"progress"`
	FinishTime  int    `json:"finishTime"`
}

type onlineStatus struct {
	Status int     `json:"Status"`
	Time   string  `json:"Time"`	
}

type Device struct {
	DeviceID           string         `json:"deviceId"`
	Connection         bool           `json:"connection"`
	InnerIP            string         `json:"innerIP"`
	DeviceName         string         `json:"deviceName"`
	Groups             []string       `json:"groups"`
	DetectionMode      bool           `json:"detectionMode"`
	ScanSchedule       []string       `json:"scanSchedule"`
	ScanFinishTime     processing     `json:"scanFinishTime"`
	CollectSchedule    dateForm       `json:"collectSchedule"`
	CollectFinishTime  processing     `json:"collectFinishTime"`
	FileSchedule       dateForm       `json:"fileDownloadDate"`
	FileFinishTime     processing     `json:"fileFinishTime"`
	ImageFinishTime    processing     `json:"imageFinishTime"`
}

func ProcessDeviceData() ([]Device, error) {
	
	devices := []Device{}
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

        var d Device
		var process, network int
		var scanSchedule, collectSchedule, fileSchedule sql.NullString
		var scanFinishTime, collectFinishTime, fileFinishTime, imageFinishTime sql.NullString
		initialProcessing := processing {
			IsFinish: true,
			Progress: 100,
			FinishTime: 1689120949,
		}
		d.ScanFinishTime = initialProcessing
		d.CollectFinishTime = initialProcessing
		d.FileFinishTime = initialProcessing
		d.ImageFinishTime = initialProcessing

        err = rows.Scan(
            &d.DeviceID,
            &d.InnerIP,
			&network,
			&process,
            &d.DeviceName,
			&scanSchedule,
			&scanFinishTime,
			&collectSchedule,
			&collectFinishTime,
			&fileSchedule,
			&fileFinishTime,
			&imageFinishTime,
        )
        if err != nil {
			return devices, err
        }

		// process detection mode
		if process == 0 && network == 0 {
			d.DetectionMode = false
		} else {
			d.DetectionMode = true
		}

		// process connection
		var status onlineStatus
		statusString, err := redis.Redis_get(d.DeviceID)
		if err != nil {
			return devices, err
		}
		err = json.Unmarshal([]byte(statusString), &status)
		if err != nil {
			return devices, err
		}
		if status.Status == 1 {
			d.Connection = true
		} else {
			d.Connection = false
		}

		// process schedule
		d.ScanSchedule = processScanSchedule(scanSchedule)
		d.CollectSchedule, err = processSchedule(collectSchedule)
		if err != nil {
			return devices, err
		}
		d.FileSchedule, err = processSchedule(fileSchedule)
		if err != nil {
			return devices, err
		}
		
		// process ScanFinishTime

		// process CollectFinishTime

		// process FileFinishTime

		// process ImageFinishTime
		d.ScanSchedule = []string {"11", "13","20"}
		d.CollectSchedule = dateForm{Date:"10", Time:"10"}
		d.FileSchedule = dateForm{Date:"10", Time:"10"}
        devices = append(devices, d)
	}

	return devices, nil
}

func processSchedule(schedule sql.NullString) (dateForm, error) {
	output := dateForm{Date: "", Time: "",}
	// handle NULL
	if !schedule.Valid {
		return output, nil
	}
	// schedule: date|time e.g. 11|20
	partitions := strings.Split(schedule.String, "|")
	if len(partitions) == 2 {
		output = dateForm {
			Date: partitions[0],
			Time: partitions[1],
		}
	} else {
		return output, fmt.Errorf("invalid schedule format")
	}
	return output, nil
}

func processScanSchedule(schedule sql.NullString) []string {
	// handle NULL
	if !schedule.Valid {
		return nil
	}
	// schedule: time,time,time e.g. 11,13,18
	partitions := strings.Split(schedule.String, ",")
	return partitions
}

func ProcessFinishTime() {

}

func CheckTaskStatus(deviceId string, work string) error {

	var status int
	query := "SELECT status FROM task WHERE client_id = ? AND type = ? AND status != 3"
	err := mariadb.DB.QueryRow(query, deviceId, work).Scan(&status)
	if err != nil {
		return err
	}

	return nil
}