package group

import (
	"edetector_API/internal/errhandler"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb/query"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	RangeBegin  string  `json:"rangeBegin"`
	RangeEnd	string  `json:"rangeEnd"`
}

type addResponse struct {
	IsSuccess bool `json:"isSuccess"`
	GroupID   int  `json:"groupID"`
}

type groupInfoResponse struct {
	IsSuccess   bool   `json:"isSuccess"`
	GroupID     int    `json:"groupID"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RangeBegin  string `json:"rangeBegin"`
	RangeEnd    string `json:"rangeEnd"`
}

func Add(c *gin.Context) {
	req := addRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		errhandler.Handler(c, err, "Invalid request format")
		return
	}
	logger.Info("Request content: " + fmt.Sprintf("%+v", req))
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
	res := addResponse{
		IsSuccess: true,
		GroupID:   id,
	}
	c.JSON(http.StatusOK, res)
}

func Remove(c *gin.Context) {
	id, err :=  strconv.Atoi(c.Query("id"))
	if err != nil {
		errhandler.Handler(c, err, "Error converting group id")
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

func GetInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		errhandler.Handler(c, err, "Error converting group id")
		return
	}
	info, err := query.LoadGroupInfo(id)
	if err != nil {
		errhandler.Handler(c, err, "Error loading group info")
		return
	}
	res := groupInfoResponse{
		IsSuccess:   true,
		GroupID:     id,
		Name:        info.Name,
		Description: info.Description,
		RangeBegin:  info.RangeBegin,
		RangeEnd:    info.RangeEnd,
	}
	c.JSON(http.StatusOK, res)
}

func Join(c *gin.Context) {
	
}

func Leave(c *gin.Context) {
	
}