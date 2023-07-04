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
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func API_init() {

	// Load configuration
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}

	// Init Logger
	logger.InitLogger(config.Viper.GetString("WORKER_LOG_FILE"))
	fmt.Println("logger is enabled please check all out info in log file: ", config.Viper.GetString("WORKER_LOG_FILE"))

	// Connect to MariaDB
	if err := mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
	} else {
		fmt.Println("Connected to MariaDB")
	}

	// Routing
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.RedirectFixedPath = true
	router.Use(cors.New(corsConfig))

	// Functions
	router.POST("/member/login", member.Login)
	router.POST("/member/signup", member.Signup)
	router.GET("/dashboard/serverState", dashboard.ServerState)
	router.GET("/dashboard/agentState", dashboard.AgentState)
	router.GET("/dashboard/ccConnectCount", dashboard.ConnectCount)
	router.GET("/dashboard/riskProgram", dashboard.RiskProgram)
	router.GET("/dashboard/riskComputer", dashboard.RiskComputer)
	router.GET("/detect/timeList", detect.TimeList)
	router.GET("/searchEvidence/DetectDevices", searchEvidence.DetectDevices)
	router.POST("/task/changeDetectMode", task.ChangeDetectMode)

	router.Run(":5000")
}