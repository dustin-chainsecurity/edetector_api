package searchEvidence

import (
	"edetector_API/internal/Error"
	"edetector_API/pkg/mariadb"
	. "edetector_API/internal/device"
	"edetector_API/pkg/redis"
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
)

type detectDevicesResponse struct {
	IsSuccess    bool       `json:"isSuccess"`
	Data         []Device   `json:"data"`
}

func DetectDevices(c *gin.Context) {

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
		Error.Handler(c, err, "Error retrieving device info")
		return
	}
	defer rows.Close()

	devices := []Device{}

	for rows.Next() {

        var d Device
		var process, network int
        err = rows.Scan(
            &d.DeviceID,
            &d.InnerIP,
			&network,
			&process,
            &d.DeviceName,
        )
        if err != nil {
			Error.Handler(c, err, "Error scanning device info")
			return
        }

		// detection mode
		if process == 0 && network == 0 {
			d.DetectionMode = false
		} else {
			d.DetectionMode = true
		}

		// connection
		var status OnlineStatus
		statusString, err := redis.Redis_get(d.DeviceID)
		if err != nil {
			Error.Handler(c, err, "Error getting online status from redis")
			return
		}		
		err = json.Unmarshal([]byte(statusString), &status)
		if err != nil {
			Error.Handler(c, err, "Error unmarshal online status")
			return
		}
		if status.Status == 1 {
			d.Connection = true
		} else {
			d.Connection = false
		}

        devices = append(devices, d)
	}

	// Create the response object
	res := detectDevicesResponse{
		IsSuccess:    true,
		Data:         devices,
	}

	c.JSON(http.StatusOK, res)
}