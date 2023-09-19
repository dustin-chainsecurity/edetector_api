package setting

import (
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetKeyImage(c *gin.Context) {
	keyType := c.Param("type")
	keyimages, err := query.GetKeyImageByType(keyType)
	if err != nil {
		errhandler.Error(c, err, "Error retrieving key images")
		return
	}
	c.JSON(200, gin.H{	
		"isSuccess": true,
		"list": keyimages,
	})
}

func AddKeyImage(c *gin.Context) {
	var req query.KeyImage
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	if req.Type == "" {
		errhandler.Info(c, fmt.Errorf("type cannot be empty"), "Error checking type existence")
		return
	}
	if req.Path == "" && req.Keyword == "" {
		errhandler.Info(c, fmt.Errorf("path or keyword cannot be empty at the same time"), "Error checking path or keyword existence")
		return
	}
	if exist, err := query.CheckKeyImage(req.Type, req.Path, req.Keyword); err != nil {
		errhandler.Error(c, err, "Error checking key image")
		return
	} else if exist {
		errhandler.Info(c, fmt.Errorf("key image already exist"), "Error checking key image")
		return
	}
	keyImageId, err := query.AddKeyImage(req)
	if err != nil {
		errhandler.Error(c, err, "Error adding key image")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"id": keyImageId,
		"message": "Key image added",
	})
}

func DeleteKeyImage(c *gin.Context) {
	var req struct {
		Ids []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	for _, id := range req.Ids {
		if exist, err := query.CheckImageID(id); err != nil {
			errhandler.Error(c, err, "Error checking key image")
			return
		} else if !exist {
			errhandler.Error(c, fmt.Errorf("key image does not exist"), "Error checking key image")
			return
		}
	}
	if err := query.DeleteKeyImage(req.Ids); err != nil {
		errhandler.Error(c, err, "Error deleting key image")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message": "Key image deleted",
	})
}

func UpdateKeyImage(c *gin.Context) {
	var req query.KeyImage
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	if exist, err := query.CheckImageID(req.Id); err != nil {
		errhandler.Error(c, err, "Error checking key image ID")
		return
	} else if !exist {
		errhandler.Error(c, fmt.Errorf("key image does not exist"), "Error checking key image ID")
		return
	}
	if req.Type == "" {
		errhandler.Info(c, fmt.Errorf("type cannot be empty"), "Error checking type existence")
		return
	}
	if req.Path == "" && req.Keyword == "" {
		errhandler.Info(c, fmt.Errorf("path or keyword cannot be empty at the same time"), "Error checking path or keyword existence")
		return
	}
	if exist, err := query.CheckKeyImage(req.Type, req.Path, req.Keyword); err != nil {
		errhandler.Error(c, err, "Error checking key image")
		return
	} else if exist {
		errhandler.Info(c, fmt.Errorf("key image already exist"), "Error checking key image")
		return
	}
	if err := query.UpdateKeyImage(req); err != nil {
		errhandler.Error(c, err, "Error updating key image")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message": "Key image updated",
	})
}

func ResetCustomizedImage(c *gin.Context) {
	if err := query.ResetCustomizedImage(); err != nil {
		errhandler.Error(c, err, "Error resetting customized image")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message": "Customized image reset",
	})
}