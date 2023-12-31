package device

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
	Status      int    `json:"status"`
	Progress    int    `json:"progress"`
	FinishTime  int    `json:"finishTime"`
	Message     string `json:"message"`
}

// status
// 0: not started
// 1: pending
// 2: working
// 3: finished
// 4: terminated
// 5: error

type Device struct {
	DeviceID           string         `json:"deviceId"`
	Connection         bool           `json:"connection"`
	InnerIP            string         `json:"innerIP"`
	DeviceName         string         `json:"deviceName"`
	Groups             []string       `json:"groups"`
	DetectionMode      bool           `json:"detectionMode"`
	Version			   string         `json:"version"`
	ScanSchedule       []string       `json:"scanSchedule"`
	ScanFinishTime     processing     `json:"scanFinishTime"`
	CollectSchedule    dateForm       `json:"collectSchedule"`
	CollectFinishTime  processing     `json:"collectFinishTime"`
	FileSchedule       dateForm       `json:"fileDownloadDate"`
	FileFinishTime     processing     `json:"fileFinishTime"`
	ImageFinishTime    processing     `json:"imageFinishTime"`
}

type Detail struct {
	DeviceID           string         `json:"id"`
	InnerIP            string         `json:"ip"`
	Mac                string         `json:"macAddress"`
	DeviceName         string         `json:"name"`
	Version			   string         `json:"version"`
	Groups             []string       `json:"group"`
}

func ProcessDeviceDetail(r mq.RawDevice) (Detail, error) {
	var d Detail
	var err error
	d.DeviceID = r.DeviceID
	d.InnerIP = r.InnerIP
	d.Mac = r.Mac
	d.DeviceName = r.DeviceName
	d.Version = r.Version

	// group
	d.Groups, err = mq.LoadGroups(d.DeviceID)
	if err != nil {
		return d, err
	}
	return d, nil
}

func ProcessRawDevice(r mq.RawDevice) (Device, error) {

	var d Device
	var err error
	d.DeviceID = r.DeviceID
	d.InnerIP = r.InnerIP
	d.DeviceName = r.DeviceName
	d.Version = r.Version

	// group
	d.Groups, err = mq.LoadGroups(d.DeviceID)
	if err != nil {
		return d, err
	}

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
		Status:     0,
		Progress:   0,
		FinishTime: 0,
		Message:    "nil",
	}

	status, progress, err := mq.LoadTaskStatus(deviceId, work)
	if err != nil {
		return output, err
	}

	switch status {
	case -1:                                   //* no tasks yet
		return output, nil
	case 0, 1:                                 //* task not started
		output.Status = 1
		return output, nil
	case 2:                                    //* working
		output.Status = 2
		output.Progress = progress
		return output, nil
	case 3, 4:                                 //* task finished or terminated
		output.Status = status
		if finishtime.Valid {
			output.FinishTime, err = parseTimestamp(finishtime.String)
			if err != nil {
				return output, fmt.Errorf("error parsing finishtime timestamp")
			}
		} else {
			// return output, nil
			return output, fmt.Errorf(work + " for " + deviceId + " finished or terminated but NULL finish time")
		}
	case 5:                                    //* error
		output.Status = 5
		output.Message = "error occurred, task failed"
		if finishtime.Valid {
			output.FinishTime, err = parseTimestamp(finishtime.String)
			if err != nil {
				return output, fmt.Errorf("error parsing finishtime timestamp")
			}
		} else {
			return output, nil
			//! TBC: return output, fmt.Errorf(work + " for " + deviceId + " failed but NULL finish time")
		}
	}
	return output, nil
}

func CheckID(id string) error {
	if exist, err := mq.CheckDevice(id); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("deviceID " + id + " does not exist")
	}
	return nil
}

func CheckAllID(devices []string) error {
	for _, id := range devices {
		if err := CheckID(id); err != nil {
			return err
		}
	}
	return nil
}