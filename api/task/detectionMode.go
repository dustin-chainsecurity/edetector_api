package task

import (
	"edetector_API/internal/Error"
	"edetector_API/pkg/mariadb"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DetectionModeRequest struct {
	Mode        bool       `json:"mode"`
	Devices    []string    `json:"deviceId"`
}

func DetectionMode(c *gin.Context) {

	var req DetectionModeRequest
	var message = "Success"
	if err := c.ShouldBindJSON(&req); err != nil {
		Error.Handler(c, err, "Invalid request format")
		return
	}

	for _, deviceId := range req.Devices {

		// fetch agent current setting
		var process, network, m int
		if req.Mode {
			m = 1
		} else {
			m = 0
		}

		query := "SELECT processreport, networkreport FROM client_setting WHERE client_id = ?"
		err := mariadb.DB.QueryRow(query, deviceId).Scan(&process, &network)
		if err != nil {
			Error.Handler(c, err, "Error retrieving current client setting")
			return
		}

		// check if the changes are needed
		if m != process || m != network {
			_, err := AddTask(deviceId, "ChangeDetectMode", fmt.Sprintf("%d|%d", m, m))
			if err != nil {
				Error.Handler(c, err, "Error adding ChangeDetectMode task")
				return
			}
			err = updateDetectMode(m, deviceId)
			if err != nil {
				Error.Handler(c, err, "Error updating client_setting table")
				return
			}
		} else {
			message = "No changes needed"
		}
	}

	res := TaskResponse {
		IsSuccess: true,
		Message: message,
	}
	c.JSON(http.StatusOK, res)
}

func updateDetectMode(mode int, deviceId string) error {
	query := "UPDATE client_setting SET processreport = ?, networkreport = ? WHERE client_id = ?"
	_, err := mariadb.DB.Exec(query, mode, mode, deviceId)
	return err
}