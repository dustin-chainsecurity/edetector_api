package searchEvidence

import (
	"database/sql"
	"edetector_API/internal/channel"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"encoding/json"
	
	"github.com/gin-gonic/gin"
)

func RefreshDevices(c *gin.Context) {
	started_tasks := <- channel.TaskChangeChannel
	finished_tasks := <- channel.TaskChangeChannel
}

func refreshStartedTasks(started [][]string) ([]device, error) {
	devices := []device{}
	var (
		d device
		process, network int
		scanSchedule, collectSchedule, fileSchedule sql.NullString
		scanFinishTime, collectFinishTime, fileFinishTime, imageFinishTime sql.NullString
	)

	initialProcessing := processing {
		IsFinish: false,
		Progress: 10,
		FinishTime: 0,
	}
	d.ScanFinishTime = initialProcessing
	d.CollectFinishTime = initialProcessing
	d.FileFinishTime = initialProcessing
	d.ImageFinishTime = initialProcessing

	for _, task := range started {
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
		err := mariadb.DB.QueryRow(query, task[0]).Scan(
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
		d.ScanFinishTime, err = processFinishTime(d.DeviceID, "StartScan", scanFinishTime)
		if err != nil {
			return devices, err
		}
		// process CollectFinishTime
		d.CollectFinishTime, err = processFinishTime(d.DeviceID, "StartCollect", collectFinishTime)
		if err != nil {
			return devices, err
		}
		// process FileFinishTime
		d.FileFinishTime, err = processFinishTime(d.DeviceID, "StartGetDrive", fileFinishTime)
		if err != nil {
			return devices, err
		}
		// process ImageFinishTime
		d.ImageFinishTime, err = processFinishTime(d.DeviceID, "StartGetImage", imageFinishTime)
		if err != nil {
			return devices, err
		}
        devices = append(devices, d)
	}

	return devices, nil
}
