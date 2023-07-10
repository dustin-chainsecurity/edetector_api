package api

import (
	"edetector_API/config"
	"edetector_API/internal/api/task"
	"edetector_API/internal/api/member"
	"edetector_API/internal/api/dashboard"
	"edetector_API/internal/api/detect"
	"edetector_API/internal/api/searchEvidence"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"fmt"
	"net"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var err error
var conn net.Conn

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
	router.POST("/searchEvidence/sendMission", searchEvidence.SendMission)
    router.GET("/searchEvidence/refreshDevices", func(c *gin.Context) {
        webSocket(c.Writer, c.Request)
    })	

	// Dashboard
	router.GET("/dashboard/serverState", dashboard.ServerState)
	router.GET("/dashboard/agentState", dashboard.AgentState)
	router.GET("/dashboard/ccConnectCount", dashboard.ConnectCount)
	router.GET("/dashboard/riskProgram", dashboard.RiskProgram)
	router.GET("/dashboard/riskComputer", dashboard.RiskComputer)
	router.GET("/detect/timeList", detect.TimeList)

	// User Tasks
	taskGroup := router.Group("/task")
	taskGroup.Use(EstablishTCPConnection()) 
	{
		taskGroup.POST("/changeDetectMode", task.ChangeDetectMode)
		taskGroup.POST("/startScan", task.StartScan)
		taskGroup.POST("/startGetDrive", task.StartGetDrive)
		taskGroup.POST("/startCollect", task.StartCollect)
	}

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