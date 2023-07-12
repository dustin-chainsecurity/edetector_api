package api

import (
	"edetector_API/config"
	"edetector_API/api/task"
	"edetector_API/api/member"
	"edetector_API/api/dashboard"
	"edetector_API/api/searchEvidence"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var err error

func Main() {

	API_init()

	// Routing
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.RedirectFixedPath = true
	router.Use(cors.New(corsConfig))

	// Login
	router.POST("/member/login", member.Login)
	router.POST("/member/loginWithToken", member.LoginWithToken)

	// Search Evidence
	router.GET("/searchEvidence/detectDevices", searchEvidence.DetectDevices)
	router.GET("/searchEvidence/refreshDevices", searchEvidence.RefreshDevices)
    router.GET("/ws", func(c *gin.Context) {
        webSocket(c.Writer, c.Request)
    })

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

	router.Run(":5000")
}

func API_init() {

	// Load configuration
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}

	// Init Logger
	logger.InitLogger(config.Viper.GetString("WORKER_LOG_FILE"))
	fmt.Println("logger is enabled please check all out info in log file: ", config.Viper.GetString("WORKER_LOG_FILE"))

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