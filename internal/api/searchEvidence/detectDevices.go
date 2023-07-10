package searchEvidence

import (
	"edetector_API/pkg/logger"
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
	ScanTime           []string       `json:"scanTime"`
	ScanFinishTime     processing     `json:"scanFinishTime"`
	CertificationDate  dateForm       `json:"certificationDate"`
	TraceFinishTime    processing     `json:"traceFinishTime"`
	TableDownloadDate  dateForm       `json:"tableDownloadDate"`
	TableFinishTime    processing     `json:"tableFinishTime"`
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
		logger.Error("Error retrieving device info: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error retrieving device info": err.Error()})
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
			logger.Error("Error scanning device info: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error scanning device info": err.Error()})
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
		statusString := redis.Redis_get(d.DeviceID)
		err = json.Unmarshal([]byte(statusString), &status)
		if err != nil {
			logger.Error("Error unmarshal online status: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error unmarshal online status": err.Error()})
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