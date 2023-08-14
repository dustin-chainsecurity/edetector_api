package websocket

import (
	"context"
	"edetector_API/api/task"
	"edetector_API/config"
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

func websocket_init(LOG_PATH string) {
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

func Main() {
	websocket_init("WS_LOG_FILE")
	ctx, cancel := context.WithCancel(context.Background())
	Quit := make(chan os.Signal, 1)

	// Gin Settings
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create(config.Viper.GetString("WS_GIN_LOG"))
	gin.DefaultWriter = io.MultiWriter(f)

	// routing
	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Accept", "Content-Length", "Authorization", "Origin", "X-Requested-With"}
	router.RedirectFixedPath = true
	router.Use(cors.New(corsConfig))

	// websocket
	router.POST("/updateTask", task.UpdateTask)
	router.GET("/ws", func(c *gin.Context) {
		webSocket(ctx, c.Writer, c.Request)
	})

	// task scheduling
	taskGroup := router.Group("/task")
	taskGroup.Use(token.TokenAuth())
	taskGroup.POST("/scheduledScan", task.ScheduledScan)
	taskGroup.POST("/scheduledCollect", task.ScheduledCollect)
	taskGroup.POST("/scheduledDownload", task.ScheduledDownload)

	// Start websocket service
	srv := &http.Server{
		Addr:    ":" + os.Args[1],
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting websocket server: " + err.Error())
			os.Exit(1)
		}
	}()

	// Maintaining Connection with MariaDB
	go query.CheckConnection(ctx)

	// shutdown process
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	logger.Info("Shutting down Websocket server...")
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Error shutting down Websocket server: " + err.Error())
		os.Exit(1)
	}
	mariadb.DB.Close()
	redis.Redis_close()
	logger.Info("Websocket server exited")
}
