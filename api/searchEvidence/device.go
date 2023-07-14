package searchEvidence

import (
	"database/sql"
	mq "edetector_API/pkg/mariadb/query"
	rq "edetector_API/pkg/redis/query"
	"fmt"
	"strings"
	"time"

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

type device struct {
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


func processRawDevice(r mq.RawDevice) (device, error) {

	var d device
	var err error
	d.DeviceID = r.DeviceID
	d.InnerIP = r.InnerIP
	d.DeviceName = r.DeviceName

	// detection mode
	if r.Network == 0 && r.Process == 0 {
		d.DetectionMode = false
	} else {
		d.DetectionMode = true
	}

	// connection
	d.Connection, err = rq.LoadOnlineStatus(d.DeviceID)
	if err != nil {
		return d, err
	}

	// schedule
	d.ScanSchedule = processScanSchedule(r.ScanSchedule)
	d.CollectSchedule, err = processSchedule(r.CollectSchedule)
	if err != nil {
		return d, err
	}
	d.FileSchedule, err = processSchedule(r.FileSchedule)
	if err != nil {
		return d, err
	}

	// finish time
	d.ScanFinishTime, err = processFinishTime(d.DeviceID, "StartScan", r.ScanFinishTime)
	if err != nil {
		return d, err
	}
	d.CollectFinishTime, err = processFinishTime(d.DeviceID, "StartCollect", r.CollectFinishTime)
	if err != nil {
		return d, err
	}
	d.FileFinishTime, err = processFinishTime(d.DeviceID, "StartGetDrive", r.FileFinishTime)
	if err != nil {
		return d, err
	}
	d.ImageFinishTime, err = processFinishTime(d.DeviceID, "StartGetImage", r.ImageFinishTime)
	if err != nil {
		return d, err
	}

	return d, nil
}

func processSchedule(schedule sql.NullString) (dateForm, error) {
	// handle NULL
	if !schedule.Valid {
		return dateForm{}, nil
	}
	// schedule: date|time e.g. 11|20
	partitions := strings.Split(schedule.String, "|")
	if len(partitions) != 2 {
		return dateForm{}, fmt.Errorf("invalid schedule format")
	}
	return dateForm {
		Date: partitions[0],
		Time: partitions[1],
	}, nil
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

func parseTimestamp(timestamp string) (int, error) {
	layout := "2006-01-02 15:04:05"
	parsedTimestamp, err := time.Parse(layout, timestamp)
	if err != nil {
		return -1, err
	}
	return int(parsedTimestamp.Unix()), nil
}

func processFinishTime(deviceId string, work string, finishtime sql.NullString) (processing, error) {

	output := processing {
		IsFinish: true,
		Progress: 0,
		FinishTime: 0,
	}

	status, progress, err := mq.LoadTaskStatus(deviceId, work)
	if err != nil {
		return output, err
	}
	switch status {
	case -1:                                   //* no tasks yet
		return output, nil
	case 0, 1:                                 //* tasks not started
		output.IsFinish = false
		return output, nil
	case 2:                                    //* working
		output.IsFinish = false
		output.Progress = progress //! TO BE ADDED: progress function
		return output, nil
	case 3:                                    //* all tasks finished
		if finishtime.Valid {
			output.FinishTime, err = parseTimestamp(finishtime.String)
			if err != nil {
				return output, fmt.Errorf("error parsing finishtime timestamp")
			}
		} else {
			return output, nil
			//! TO BE CHANGED: return output, fmt.Errorf(work + " for " + deviceId + " finished but NULL finish time")
		}
	}
	return output, nil
}