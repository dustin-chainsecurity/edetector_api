package taskservice

import (
	"context"
	"edetector_API/config"
	"edetector_API/internal/schedule"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var err error

func taskservice_init(LOG_PATH string, HOSTNAME string, APP string) {
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
	// Connect to MariaDB
	if err = mariadb.Connect_init(); err != nil {
		logger.Error("Error connecting to mariadb: " + err.Error())
		return
	}
}

func Start() {
	taskservice_init("TASK_LOG_FILE", "taskservice", "TASKSERVICE")
	Quit := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	logger.Info("Task service enabled...")
	go Main(ctx)
	go schedule.ScheduleTask(ctx)
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	cancel()
	mariadb.DB.Close()
	redis.Redis_close()
	logger.Info("Task service shutdown")
}

func Main(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Task service is shutting down...")
			return
		default:
			findTask(ctx)
			time.Sleep(5 * time.Second)
		}
	}
}