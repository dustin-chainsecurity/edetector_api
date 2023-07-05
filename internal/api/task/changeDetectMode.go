package task

import (
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ChangeDetectMode(c *gin.Context) {

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// fetch agent current setting
	var process, network int
	query := "SELECT processreport, networkreport FROM client_setting WHERE client_id = ?"
	err := mariadb.DB.QueryRow(query, req.Key).Scan(&process, &network)
	if err != nil {
		logger.Error("Error retrieving current client setting: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error retrieving current client setting": err.Error()})
		return
	}

	// check if the changes are needed
	values := strings.Split(req.Message, "|")  // "0|0"
	// Convert string values to integers
	new_process, err1 := strconv.Atoi(values[0])
	new_network, err2 := strconv.Atoi(values[1])
	if err1 != nil || err2 != nil {
		logger.Error("Invalid message input: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Invalid message input": err.Error()})
		return
	}
	if new_process == process && new_network == network {
		res := TaskResponse {
			IsSuccess: true,
			Message: "Nothing changed",
		}
		c.JSON(http.StatusOK, res)
	}

	// update agent settings
	query = "UPDATE user_info SET processreport = ?, networkreport = ? WHERE id = ?"
	_, err = mariadb.DB.Exec(query, new_process, new_network)
	if err != nil {
		logger.Error("Error updating client setting: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error updating client setting": err.Error()})
		return
	}

	response, err := SendTCPRequest(c, "ChangeDetectMode", req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		logger.Error("Error sending TCP request: " + err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}