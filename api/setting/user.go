package setting

import (
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	users, err := query.GetUsers()
	if err != nil {
		errhandler.Error(c, err, "Error retrieving users")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"users": users,
	})
}

func AddUser(c *gin.Context) {
	var user query.User
	if err := c.ShouldBindJSON(&user); err != nil {	
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", user))
	if user.Username == "" || user.Password == "" {
		errhandler.Info(c, nil, "Username or password cannot be empty")
		return
	}
	if exist, err := query.CheckUsername(user.Username); err != nil {
		errhandler.Error(c, err, "Error checking username existence")
		return
	} else if exist {
		errhandler.Info(c, err, "Username already exist")
		return
	}
	userId, err := query.AddUser(user); 
	if err != nil {
		errhandler.Error(c, err, "Error adding user")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"userID": userId,
		"message": "User add success",
	})
}

func UpdateUserInfo(c *gin.Context) {
	var user query.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", user))
	if user.Username == "" || user.Password == "" {
		errhandler.Info(c, fmt.Errorf("username or password cannot be empty"), "Error checking username existence")
		return
	}
	if exist, err := query.CheckUserId(user.UserID); err != nil {
		errhandler.Error(c, err, "Error checking user existence")
		return
	} else if !exist {
		errhandler.Error(c, fmt.Errorf("userID does not exist"), "Error checking user existence")
		return
	}
	if err := query.UpdateUser(user); err != nil {
		errhandler.Error(c, err, "Error updating user")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message": "User update success",
	})
}

func DeleteUser(c *gin.Context) {
	var req struct {
		IDS []int `json:"ids"`
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Info(c, err, "Invalid request format")
		return
	}
	if err := query.DeleteUser(req.IDS); err != nil {
		errhandler.Error(c, err, "Error deleting user")
		return
	}
	c.JSON(200, gin.H{
		"isSuccess": true,
		"message": "User delete success",
	})
}

