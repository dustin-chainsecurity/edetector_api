package updatedb

import (
	"crypto/rand"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"encoding/json"
	"fmt"
)

type ClientOnlineStatus struct {
	Status int
	Time   string
}

func InsertData() {

	var key, ip, deviceName, query, mac string
	var err error

	for i := 1; i < 6; i++ {

		key = fmt.Sprintf("test_agent%d", i)
		ip = fmt.Sprintf("192.168.0.%d", i)
		mac = generateMACAddress()
		deviceName = fmt.Sprintf("PC-%d", i)

		// client
		query = "INSERT INTO client (client_id, ip, mac) VALUES (?, ?, ?)"
		_, err = mariadb.DB.Exec(query, key, ip, mac)
		if err != nil {
			logger.Error("Error storing client data: " + err.Error())
			return
		}

		// client_setting
		query = "INSERT INTO client_setting (client_id, networkreport, processreport) VALUES (?, ?, ?)"
		_, err = mariadb.DB.Exec(query, key, 0, 0)
		if err != nil {
			logger.Error("Error storing client data: " + err.Error())
			return
		}

		// client_info
		query = "INSERT INTO client_info (client_id, sysinfo, osinfo, computername, username, fileversion, boottime) VALUES (?, ?, ?, ?, ?, ?, ?)"
		_, err = mariadb.DB.Exec(query, key, "x64", "Windows 10 Home", deviceName, "SYSTEM", "3.4.2.0,1988,1989", "20230706110846")
		if err != nil {
			logger.Error("Error storing client data: " + err.Error())
			return
		}

		// client_task_status
		query = "INSERT INTO client_task_status (client_id) VALUES (?)"
		_, err = mariadb.DB.Exec(query, key)
		if err != nil {
			logger.Error("Error storing client_task_status data: " + err.Error())
			return
		}

		onlineStatusInfo := ClientOnlineStatus {
			Status: 0,
			Time:   "20230706110846",
		}
		redis.Redis_set(key, onlineStatusInfo.Marshal())
	}
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