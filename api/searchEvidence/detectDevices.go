package searchEvidence

import (
	"edetector_API/internal/Error"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
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
	SubGroup           []string       `json:"subGroup"`
	DetectionMode      bool           `json:"detectionMode"`
	ScanSchedule       []string       `json:"scanTime"`
	ScanFinishTime     processing     `json:"scanFinishTime"`
	CollectSchedule    dateForm       `json:"certificationDate"`
	CollectFinishTime  processing     `json:"traceFinishTime"`
	FileSchedule       dateForm       `json:"tableDownloadDate"`
	FileFinishTime     processing     `json:"tableFinishTime"`
	ImageFinishTime    processing     `json:"imageFinishTime"`
}

type detectDevicesResponse struct {
	IsSuccess    bool      `json:"isSuccess"`
	Data         []device  `json:"data"`
}

type onlineStatus struct {
	Status int     `json:"Status"`
	Time   string  `json:"Time"`	
}

func DetectDevices(c *gin.Context) {

	query := `
		SELECT C.client_id, C.ip, S.networkreport, S.processreport, I.computername 
		FROM client AS C
		JOIN client_setting AS S ON C.client_id = S.client_id
		JOIN client_info AS I ON S.client_id = I.client_id
	`
	rows, err := mariadb.DB.Query(query)
	if err != nil {
		Error.Handler(c, err, "Error retrieving device info")
		return
	}
	defer rows.Close()

	devices := []device{}

	for rows.Next() {

        var d device
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
		var status onlineStatus
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