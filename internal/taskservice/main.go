package taskservice

import (
	"context"
	"edetector_API/api"
	"edetector_API/internal/schedule"
	"edetector_API/pkg/logger"
	"edetector_API/pkg/mariadb"
	"edetector_API/pkg/redis"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start() {
	api.API_init("TASK_LOG_FILE")
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