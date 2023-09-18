package api

import (
	"context"
	"edetector_API/api/analysis"
	"edetector_API/api/clear"
	"edetector_API/api/group"
	"edetector_API/api/member"
	"edetector_API/api/saveagent"
	"edetector_API/api/searchEvidence"
	"edetector_API/api/setting"
	"edetector_API/api/task"
	"edetector_API/api/testing"
	"edetector_API/config"
	"edetector_API/pkg/fflag"
	"edetector_API/internal/token"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/mariadb/query"
	"edetector_API/pkg/redis"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var err error

func API_init(LOG_PATH string, HOSTNAME string, APP string) {
	// Load feature flag
	fflag.Get_fflag()
	if fflag.FFLAG == nil {
		logger.Error("Error loading feature flag")
		return
	}	
	// Load configuration
	if config.LoadConfig() == nil {
		fmt.Println("Error loading config file")
		return
	}
	// Init Logger
	logger.InitLogger(config.Viper.GetString(LOG_PATH), HOSTNAME, APP)
	logger.Log.Info("Logger enabled, log file: " + config.Viper.GetString(LOG_PATH))
	// Connect to Redis
	if db := redis.Redis_init(); db == nil {
		logger.Error("Error connecting to redis")
		return
	}
	logger.Info("Redis connected")
	// Connect to MariaDB
	if err = mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
		return
	}
	logger.Info("MariaDB connected")
}

func Main(version string) {
	API_init("API_LOG_FILE", "api", "API")
	logger.Info("API server started, version: " + version)
	ctx, cancel := context.WithCancel(context.Background())
	Quit := make(chan os.Signal, 1)

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
	router.Use(logger.GinLog())

	// Download Executables
	router.StaticFS("/server_file", http.Dir("./api/server_files"))

	// Backend
	router.GET("/save", saveagent.SaveAgent)
	router.POST("/sendMission", task.SendMission)
	router.GET("/test", testing.Test)
	router.POST("/updateProgress", testing.UpdateProgress)
	router.POST("/addDevice", testing.AddDevice)

	// Admin Group
	adminGroup := router.Group("/admin")
	adminGroup.Use(token.TokenAdminAuth())
	adminGroup.DELETE("/rabbit", clear.ClearRabbit)
	adminGroup.DELETE("/redis", clear.ClearRedis)
	adminGroup.DELETE("/maria", clear.ClearMaria)
	adminGroup.DELETE("/elastic", clear.ClearElastic)
	adminGroup.POST("/signup", setting.AddUser)

	// Login
	router.POST("/login", member.Login)
	router.POST("/loginWithToken", member.LoginWithToken)

	// Member Group
	memberGroup := router.Group("/member")
	memberGroup.Use(token.TokenAuth())

	// Search Evidence Page
	searchEvidenceGroup := router.Group("/searchEvidence")
	searchEvidenceGroup.Use(token.TokenAuth())
	searchEvidenceGroup.GET("/detectDevices", searchEvidence.DetectDevices)
	searchEvidenceGroup.POST("/refresh", searchEvidence.Refresh)

	// Analysis Page
	analysisGroup := router.Group("/analysisPage")
	analysisGroup.Use(token.TokenAuth())
	analysisGroup.GET("/allDeviceDetail", analysis.DeviceDetail)
	analysisGroup.POST("/template", analysis.AddTemplate)
	analysisGroup.GET("/template", analysis.GetTemplateList)
	analysisGroup.GET("/template/:id", analysis.GetTemplate)
	analysisGroup.PUT("/template/:id", analysis.UpdateTemplate)
	analysisGroup.DELETE("/template/:id", analysis.DeleteTemplate)

	// Working Server Tasks
	taskGroup := router.Group("/task")
	taskGroup.Use(token.TokenAuth())
	taskGroup.POST("/sendMission", task.SendMission)
	taskGroup.POST("/detectionMode", task.DetectionMode)
	
	// Group
	groupGroup := router.Group("/group")
	groupGroup.Use(token.TokenAuth())
	groupGroup.POST("", group.Add)
	groupGroup.GET("", group.GetList)
	groupGroup.POST("/add", group.Add)
	groupGroup.GET("/all", group.GetList)
	groupGroup.GET("/:id", group.GetInfo)
	groupGroup.PUT("/:id", group.Update)
	groupGroup.DELETE("/:id", group.Remove)
	groupGroup.POST("/device", group.Join)
	groupGroup.DELETE("/device", group.Leave)

	// Setting Group
	settingGroup := router.Group("/setting")
	settingGroup.Use(CORS())
	settingGroup.Use(token.TokenAuth())
	settingGroup.GET("/system/:field", setting.GetSettingField)
	settingGroup.POST("/system/:field", setting.UpdateSettingField)

	settingGroup.GET("/user", setting.GetUserInfo)
	settingGroup.POST("/user", setting.AddUser)
	settingGroup.PUT("/user", setting.UpdateUserInfo)
	settingGroup.DELETE("/user", setting.DeleteUser)

	settingGroup.GET("/whitelist", setting.GetWhiteList)
	settingGroup.POST("/whitelist", setting.AddWhiteList)
	settingGroup.DELETE("/whitelist", setting.DeleteWhiteList)
	settingGroup.PUT("/whitelist", setting.UpdateWhiteList)
	settingGroup.GET("/blacklist", setting.GetBlackList)
	settingGroup.POST("/blacklist", setting.AddBlackList)
	settingGroup.DELETE("/blacklist", setting.DeleteBlackList)
	settingGroup.PUT("/blacklist", setting.UpdateBlackList)
	settingGroup.GET("/hacklist", setting.GetHackList)
	settingGroup.POST("/hacklist", setting.AddHackList)
	settingGroup.DELETE("/hacklist", setting.DeleteHackList)
	settingGroup.PUT("/hacklist", setting.UpdateHackList)

	settingGroup.GET("keyImage/:type", setting.GetKeyImage)
	settingGroup.POST("keyImage", setting.AddKeyImage)
	settingGroup.DELETE("keyImage", setting.DeleteKeyImage)
	settingGroup.PUT("keyImage", setting.UpdateKeyImage)

	settingGroup.GET("/logs", setting.GetLogs)
	settingGroup.DELETE("/removeAgent", setting.RemoveAgent)
	settingGroup.GET("/license", setting.GetLicense)

	// Start API service
	srv := &http.Server{
		Addr:    ":" + os.Args[1],
		Handler: router,
	}
	go func () {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting API server: " + err.Error())
			os.Exit(1)
		}
	}()

	// Maintaining Connection with MariaDB
	go query.CheckConnection(ctx)

	// Graceful Shutdown
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	logger.Info("Shutting down API server...")
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Error shutting down API server: " + err.Error())
		os.Exit(1)
	}
	mariadb.DB.Close()
	redis.Redis_close()
	logger.Info("API server exited")
}

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}