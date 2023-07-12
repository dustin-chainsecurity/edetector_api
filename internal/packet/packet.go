package packet

import (
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"encoding/json"
	"fmt"
	"strings"
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
		var scanSchedule, collectSchedule, fileSchedule, scanFinishTime, collectFinishTime, fileFinishTime, imageFinishTime string

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
		
		// process finish time

        devices = append(devices, d)
	}

	return devices, nil
}

func processSchedule(schedule string) (dateForm, error) {
	// schedule: date|time e.g. 11|20
	partitions := strings.Split(schedule, "|")
	output := dateForm{}
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

func processScanSchedule(schedule string) []string {
	// schedule: time,time,time e.g. 11,13,18
	partitions := strings.Split(schedule, ",")
	return partitions
}

func processFinishTime() {

}