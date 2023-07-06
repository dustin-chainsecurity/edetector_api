package main

import (
	"crypto/rand"
	"edetector_API/internal/api"
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

func main() {

	api.API_init()

	var key, ip, deviceName, query, mac string
	var err error

	for i := 200; i < 501; i++ {

		key = fmt.Sprintf("deviceIDNumber%d", i)
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

		// client_group
		query = "INSERT INTO client_group (client_id, group_id) VALUES (?, ?)"
		_, err = mariadb.DB.Exec(query, key, 1)
		if err != nil {
			logger.Error("Error storing client data: " + err.Error())
			return
		}
		query = "INSERT INTO client_group (client_id, group_id) VALUES (?, ?)"
		_, err = mariadb.DB.Exec(query, key, 2)
		if err != nil {
			logger.Error("Error storing client data: " + err.Error())
			return
		}

		onlineStatusInfo := ClientOnlineStatus {
			Status: 1,
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
	// Create a byte slice of length 6 to store the MAC address
	mac := make([]byte, 6)

	// Generate 6 random bytes
	_, err := rand.Read(mac)
	if err != nil {
		panic(err)
	}

	// Set the locally administered bit (the second least significant bit of the first byte)
	mac[0] |= 2

	// Format the MAC address as a string
	macAddress := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])

	return macAddress
}