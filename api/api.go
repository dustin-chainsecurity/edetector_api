package api

import (
	"edetector_API/api/dashboard"
	"edetector_API/api/member"
	"edetector_API/api/saveagent"
	"edetector_API/api/searchEvidence"
	"edetector_API/api/task"
	"edetector_API/config"
	"edetector_API/internal/token"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var err error

func Main() {

	API_init()
	Quit := make(chan os.Signal, 1)

	// Routing
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Accept", "Content-Length", "Authorization", "Origin", "X-Requested-With"}
	router.RedirectFixedPath = true
	router.Use(cors.New(corsConfig))

	// Backend
	router.GET("/check", task.Check)
	router.GET("/save", saveagent.SaveAgent)
	router.POST("/sendMission", task.SendMission)
	// QA
	router.POST("/updateProgress", task.UpdateProgress)

	// Login
	router.POST("/member/signup", member.Signup)
	router.POST("/member/login", member.Login)
	router.POST("/member/loginWithToken", member.LoginWithToken)

	// Use Token Authentication
	router.Use(token.TokenAuth())

	// Search Evidence
	router.GET("/searchEvidence/detectDevices", searchEvidence.DetectDevices)
	router.POST("/searchEvidence/refresh", searchEvidence.Refresh)

	// Working Server Tasks
	taskGroup := router.Group("/task")
	taskGroup.POST("/sendMission", task.SendMission)
	taskGroup.POST("/detectionMode", task.DetectionMode)

	// Dashboard
	router.GET("/dashboard/serverState", dashboard.ServerState)
	router.GET("/dashboard/agentState", dashboard.AgentState)
	router.GET("/dashboard/ccConnectCount", dashboard.ConnectCount)
	router.GET("/dashboard/riskProgram", dashboard.RiskProgram)
	router.GET("/dashboard/riskComputer", dashboard.RiskComputer)

	// shutdown process
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-Quit
		fmt.Println("Web API shutdown complete.")
		os.Exit(0)
	}()

	router.Run(":" + os.Args[1])
}

func API_init() {

	// Load configuration
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}

	// Init Logger
	logger.InitLogger(config.Viper.GetString("WORKER_LOG_FILE"))
	fmt.Println("Logger enabled, log file: ", config.Viper.GetString("WORKER_LOG_FILE"))

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