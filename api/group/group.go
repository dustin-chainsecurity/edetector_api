package group

import (
	"edetector_API/internal/device"
	"edetector_API/pkg/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type updateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RangeBegin  string `json:"rangeBegin"`
	RangeEnd    string `json:"rangeEnd"`
}

type addGroupResponse struct {
	IsSuccess bool `json:"isSuccess"`
	GroupID   int  `json:"groupID"`
}

type groupInfoResponse struct {
	IsSuccess   bool     `json:"isSuccess"`
	GroupID     int      `json:"groupID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RangeBegin  string   `json:"rangeBegin"`
	RangeEnd    string   `json:"rangeEnd"`
	Devices     []string `json:"devices"`
}

type updateDevicesRequest struct {
	Groups  []int    `json:"groups"`
	Devices []string `json:"devices"`
}

type group struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RangeBegin  string   `json:"rangeBegin"`
	RangeEnd    string   `json:"rangeEnd"`
	Devices     []string `json:"devices"`
}

func Add(c *gin.Context) {
	req := updateGroupRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))

	if exist, err := query.CheckGroupName(req.Name); err != nil {
		errhandler.Handler(c, err, "Error checking group name existence")
		return
	} else if exist {
		errhandler.Handler(c, fmt.Errorf("group name "+req.Name+" already exists"), "Error checking group name")
		return
	}

	id, err := query.AddGroup(query.GroupInfo{
		Name:        req.Name,
		Description: req.Description,
		RangeBegin:  req.RangeBegin,
		RangeEnd:    req.RangeEnd,
	})
	if err != nil {
		errhandler.Handler(c, err, "Error adding group")
		return
	}
	res := addGroupResponse{
		IsSuccess: true,
		GroupID:   id,
	}
	c.JSON(http.StatusOK, res)
}

func Update(c *gin.Context) {
	req := updateGroupRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errhandler.Handler(c, err, "Error retrieving groupID")
		return
	}
	// check if the group exists
	if err := checkID(id); err != nil {
		errhandler.Handler(c, err, "Invalid group ID")
		return
	}
	// check if the group name exists
	if err := query.CheckGroupNameForUpdate(id, req.Name); err != nil {
		errhandler.Handler(c, err, "Error checking group name")
		return
	}
	err = query.UpdateGroup(query.GroupInfo{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		RangeBegin:  req.RangeBegin,
		RangeEnd:    req.RangeEnd,
	})
	if err != nil {
		errhandler.Handler(c, err, "Error updating group")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

func Remove(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errhandler.Handler(c, err, "Error retrieving groupID")
		return
	}
	// check if the group exists
	if err := checkID(id); err != nil {
		errhandler.Handler(c, err, "Invalid group ID")
		return
	}
	err = query.RemoveGroup(id)
	if err != nil {
		errhandler.Handler(c, err, "Error removing group")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

func GetList(c *gin.Context) {
	groups, err := query.LoadAllGroupInfo()
	if err != nil {
		errhandler.Handler(c, err, "Error loading group list")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
		"groups":    groups,
	})
}

func GetInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errhandler.Handler(c, err, "Error retrieving groupID")
		return
	}
	// check if the group exists
	if err := checkID(id); err != nil {
		errhandler.Handler(c, err, "Invalid group ID")
		return
	}
	info, err := query.LoadGroupInfo(id)
	if err != nil {
		errhandler.Handler(c, err, "Error loading group info")
		return
	}
	res := groupInfoResponse{
		IsSuccess:   true,
		GroupID:     info.ID,
		Name:        info.Name,
		Description: info.Description,
		RangeBegin:  info.RangeBegin,
		RangeEnd:    info.RangeEnd,
		Devices:     info.Devices,
	}
	c.JSON(http.StatusOK, res)
}

func Join(c *gin.Context) {
	req := updateDevicesRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	// check if the all group exists
	err := checkAllID(req.Groups)
	if err != nil {
		errhandler.Handler(c, err, "Invalid group ID")
		return
	}
	// check devices
	if err := device.CheckAllID(req.Devices); err != nil {
		errhandler.Handler(c, err, "Invalid device ID")
		return
	}
	err = query.AddGroupDevices(req.Groups, req.Devices)
	if err != nil {
		errhandler.Handler(c, err, "Error adding devices to group")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

func Leave(c *gin.Context) {
	// Remove devices from group
	req := updateDevicesRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
	// check if the all group exists
	err := checkAllID(req.Groups)
	if err != nil {
		errhandler.Handler(c, err, "Invalid group ID")
		return
	}
	// check devices
	if err := device.CheckAllID(req.Devices); err != nil {
		errhandler.Handler(c, err, "Invalid device ID")
		return
	}
	err = query.RemoveGroupDevices(req.Groups, req.Devices)
	if err != nil {
		errhandler.Handler(c, err, "Error removing devices from group")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"isSuccess": true,
	})
}

func checkID(id int) error {
	if exist, err := query.CheckGroupID(id); err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("groupID " + strconv.Itoa(id) + " does not exist")
	}
	return nil
}

func checkAllID(groups []int) error {
	for _, id := range groups {
		if err := checkID(id); err != nil {
			return err
		}
	}
	return nil
}
