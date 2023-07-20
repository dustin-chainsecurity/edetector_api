package taskservice

import (
	"context"
	"edetector_API/api"
	"edetector_API/internal/schedule"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Start() {
	api.API_init()
	Quit := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println("Task service enabled...")
	go Main(ctx)
	go schedule.ScheduleTask(ctx)
	signal.Notify(Quit, syscall.SIGINT, syscall.SIGTERM)
	<-Quit
	cancel()
	fmt.Println("Server shutdown complete...")
}

func Main(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Task service is shutting down...")
			return
		default:
			findTask(ctx)
			time.Sleep(5 * time.Second)
		}
	}
}