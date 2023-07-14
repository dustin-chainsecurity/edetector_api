package query

import (
	"edetector_API/pkg/redis"
	"encoding/json"
)

type onlineStatus struct {
	Status int     `json:"Status"`
	Time   string  `json:"Time"`	
}

func LoadOnlineStatus(deviceId string) (bool, error) {
	var status onlineStatus
	statusString, err := redis.Redis_get(deviceId)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal([]byte(statusString), &status)
	if err != nil {
		return false, err
	}

	if status.Status == 1 {
		return true, nil
	}
	return false, nil
}