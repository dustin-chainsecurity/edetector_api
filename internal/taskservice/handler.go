package taskservice

import (
	"context"
	"edetector_API/internal/request"
	"edetector_API/pkg/logger"
	mq "edetector_API/pkg/mariadb/query"
	rq "edetector_API/pkg/redis/query"
)

var taskchans = make(map[string]chan string)

func findTask(ctx context.Context) {
	unhandled_tasks := loadUnhandledTasks()
	for _, task := range unhandled_tasks {
		status, err := rq.LoadOnlineStatus(task.clientId)
		if err != nil {
			logger.Error(err.Error())
		}
		if status {
			if _, ok := taskchans[task.clientId]; !ok {
				taskchans[task.clientId] = make(chan string, 1024)
				go taskHandler(ctx, taskchans[task.clientId], task.clientId)
			}
			mq.UpdateTaskStatus(task.taskId, 1)
			taskchans[task.clientId] <- task.taskId
		}
	}
}

func taskHandler(ctx context.Context, ch chan string, client string) {
	for {
		select {
		case <-ctx.Done():
			logger.Info("Task handler for " + client + " is shutting down...")
			return
		case taskId := <-ch:
			logger.Info("Task " + taskId + " handled for client " + client)
			mq.UpdateTaskStatus(taskId, 2)
			request.SendToServer(taskId)
			request.SendUpdateTaskToApi(client)
		}
	}
}