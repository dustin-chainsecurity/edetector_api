package websocket

import (
	"context"
	"edetector_API/api"
	"edetector_API/api/task"
	"edetector_API/internal/token"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Main() {
	api.API_init()
	ctx, cancel := context.WithCancel(context.Background())
	Quit := make(chan os.Signal, 1)

	// routing
	gin.SetMode(gin.ReleaseMode)
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

	// shutdown process
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-Quit
		cancel()
		fmt.Println("Websocket shutdown complete.")
		os.Exit(0)
	}()

	router.Run(":" + os.Args[1])
}
