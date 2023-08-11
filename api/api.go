package api

import (
	"context"
	"edetector_API/api/analysis"
	"edetector_API/api/clear"
	"edetector_API/api/group"
	"edetector_API/api/member"
	"edetector_API/api/saveagent"
	"edetector_API/api/searchEvidence"
	"edetector_API/api/task"
	"edetector_API/api/testing"
	"edetector_API/config"
	"edetector_API/internal/token"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/mariadb/query"
	"edetector_API/pkg/redis"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var err error

func Main() {

	API_init("API_LOG_FILE")
	Quit := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	logger.Info("Web API service enabled...")

	// Gin Settings
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create(config.Viper.GetString("API_GIN_LOG"))
	gin.DefaultWriter = io.MultiWriter(f)

	// Routing
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Accept", "Content-Length", "Authorization", "Origin", "X-Requested-With"}
	router.RedirectFixedPath = true
	router.Use(cors.New(corsConfig))

	// Backend
	router.GET("/save", saveagent.SaveAgent)
	router.POST("/sendMission", task.SendMission)
	router.GET("/test", testing.Test)
	router.DELETE("/rabbit", clear.ClearRabbit)
	router.DELETE("/redis", clear.ClearRedis)
	router.DELETE("/maria", clear.ClearMaria)
	router.DELETE("/elastic", clear.ClearElastic)

	// Testing
	router.POST("/updateProgress", testing.UpdateProgress)
	router.POST("/addDevice", testing.AddDevice)

	// Login
	router.POST("/member/login", member.Login)
	router.POST("/member/loginWithToken", member.LoginWithToken)

	// Use Token Authentication
	router.Use(token.TokenAuth())
	router.POST("/member/signup", member.Signup)

	// Search Evidence
	router.GET("/searchEvidence/detectDevices", searchEvidence.DetectDevices)
	router.POST("/searchEvidence/refresh", searchEvidence.Refresh)

	// Analysis Page
	router.GET("/analysisPage/allDeviceDetail", analysis.DeviceDetail)
	router.POST("/analysisPage/template", analysis.AddTemplate)
	router.GET("/analysisPage/template", analysis.GetTemplateList)
	router.GET("/analysisPage/template/:id", analysis.GetTemplate)
	router.PUT("/analysisPage/template/:id", analysis.UpdateTemplate)
	router.DELETE("/analysisPage/template/:id", analysis.DeleteTemplate)

	// Working Server Tasks
	taskGroup := router.Group("/task")
	taskGroup.POST("/sendMission", task.SendMission)
	taskGroup.POST("/detectionMode", task.DetectionMode)

	// Group
	router.POST("/group", group.Add)
	router.GET("/group", group.GetList)
	router.GET("/group/:id", group.GetInfo)
	router.PUT("/group/:id", group.Update)
	router.DELETE("/group/:id", group.Remove)
	router.POST("/group/device", group.Join)
	router.DELETE("/group/device", group.Leave)

	// Maintaining Connection with MariaDB
	go query.CheckConnection(ctx)

	// Shutdown Process
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-Quit
		cancel()
		logger.Info("Shutting down API server...")
		os.Exit(0)
	}()

	router.Run(":" + os.Args[1])
}

func API_init(LOG_PATH string) {
	// Load configuration
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}
	// Init Logger
	logger.InitLogger(config.Viper.GetString(LOG_PATH))
	logger.Log.Info("Logger enabled, log file: " + config.Viper.GetString(LOG_PATH))
	// Connect to Redis
	if db := redis.Redis_init(); db == nil {
		logger.Error("Error connecting to redis")
		return
	}
	// Connect to MariaDB
	if err = mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
		return
	}
}