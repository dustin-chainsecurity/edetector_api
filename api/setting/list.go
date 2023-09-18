package setting

import (
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetWhiteList(c *gin.Context) {
	whitelists, err := query.GetList("white_list")
	if err != nil {
		errhandler.Error(c, err, "Error retreiving whitelists")
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"list":   whitelists,
	})
}

func AddWhiteList(c *gin.Context) {
	var req query.List
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	id, err := query.AddList("white_list", req)
	if err != nil {
		errhandler.Error(c, err, "Error adding whitelist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"id":   id,
	})
}

func UpdateWhiteList(c *gin.Context) {
	var req query.List
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	err := query.UpdateList("white_list", req)
	if err != nil {
		errhandler.Error(c, err, "Error updating whitelist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message": "Update success",
	})	
}

func DeleteWhiteList(c *gin.Context) {
	var req struct {
		IDS []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	err := query.DeleteList("white_list", req.IDS)
	if err != nil {
		errhandler.Error(c, err, "Error deleting whitelist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
	})
}

func GetBlackList(c *gin.Context) {
	blacklists, err := query.GetList("black_list")
	if err != nil {
		errhandler.Error(c, err, "Error retrieving blacklists")
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"list":   blacklists,
	})
}

func AddBlackList(c *gin.Context) {
	var req query.List
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	id, err := query.AddList("black_list", req)
	if err != nil {
		errhandler.Error(c, err, "Error adding blacklist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"id":   id,
	})
}

func UpdateBlackList(c *gin.Context) {
	var req query.List
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	err := query.UpdateList("black_list", req)
	if err != nil {
		errhandler.Error(c, err, "Error updating blacklist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message": "Update success",
	})
}

func DeleteBlackList(c *gin.Context) {
	var req struct {
		IDS []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	err := query.DeleteList("black_list", req.IDS)
	if err != nil {
		errhandler.Error(c, err, "Error deleting blacklist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
	})
}

func GetHackList(c *gin.Context) {
	hacklists, err := query.GetHackList()
	if err != nil {
		errhandler.Error(c, err, "Error retrieving hacklists")
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"list":   hacklists,
	})
}

func AddHackList(c *gin.Context) {
	var req query.HackList
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	id, err := query.AddHackList(req)
	if err != nil {
		errhandler.Error(c, err, "Error adding hacklist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"id":   id,
	})
}

func UpdateHackList(c *gin.Context) {
	var req query.HackList
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	err := query.UpdateHackList(req)
	if err != nil {
		errhandler.Error(c, err, "Error updating hacklist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message": "Update success",
	})
}

func DeleteHackList(c *gin.Context) {
	var req struct {
		IDS []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	err := query.DeleteHackList(req.IDS)
	if err != nil {
		errhandler.Error(c, err, "Error deleting hacklist")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
	})
}