package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type certificationDate struct {
	Date int    `json:"date"`
	Time string `json:"time"`
}

type tableDownloadDate struct {
	Date int    `json:"date"`
	Time string `json:"time"`
}

type device struct {
	DeviceID           string             `json:"deviceId"`
	Connection         bool               `json:"connection"`
	InnerIP            string             `json:"innerIP"`
	DeviceName         string             `json:"deviceName"`
	SubGroup           []string           `json:"subGroup"`
	DetectionMode      bool               `json:"detectionMode"`
	ScanTime           []string           `json:"scanTime"`
	ScanFinishTime     int64              `json:"scanFinishTime"`
	CertificationDate  certificationDate  `json:"CertificationDate"`
	TraceFinishTime    int64              `json:"traceFinishTime"`
	TableDownloadDate  tableDownloadDate  `json:"tableDownloadDate"`
	TableFinishTime    int64              `json:"tableFinishTime"`
	ImageFinishTime    int64              `json:"ImageFinishTime"`
}

type detectDevicesResponse struct {
	IsSuccess    bool     `json:"isSuccess"`
	TotalPages   int      `json:"totalPages"`
	TotalDevices int      `json:"totalDevices"`
	Data         []device `json:"Data"`
}

func detectDevices(c *gin.Context) {
	// page := c.Query("pages")

	// Simulate generating the response data based on the requested page
	data := []device{
		{
			DeviceID:          "deviceIDNumber1",
			Connection:        true,
			InnerIP:           "1.2.4.69.69",
			DeviceName:        "PC-01",
			SubGroup:          []string{"group1", "group2"},
			DetectionMode:     true,
			ScanTime:          []string{"0000", "1200", "2400"},
			ScanFinishTime:    1687332360,
			CertificationDate: certificationDate{Date: 1, Time: "2400"},
			TraceFinishTime:   1687332360,
			TableDownloadDate: tableDownloadDate{Date: 1, Time: "2400"},
			TableFinishTime:   1687332360,
			ImageFinishTime:   1687332360,
		},
	}

	// Create the response object
	res := detectDevicesResponse{
		IsSuccess:    true,
		TotalPages:   5,
		TotalDevices: 49,
		Data:         data,
	}

	c.JSON(http.StatusOK, res)
}