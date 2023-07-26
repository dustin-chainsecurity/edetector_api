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
	unhandle_tasks := loadUnhandleTasks()
	for _, task := range unhandle_tasks {
		status, err := rq.LoadOnlineStatus(task.clientId)
		if err != nil {
			logger.Error(err.Error())
		}
		if status {
			if _, ok := taskchans[task.clientId]; !ok {
				taskchans[task.clientId] = make(chan string, 1024)
				go taskHandler(ctx, taskchans[task.clientId], task.clientId)
			}
			taskchans[task.clientId] <- task.taskId
			mq.UpdateTaskStatus(task.taskId, 1)
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
