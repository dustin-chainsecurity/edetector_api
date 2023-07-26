package schedule

import (
	"context"
	"edetector_API/internal/request"
	"edetector_API/pkg/logger"
	"fmt"
	"strconv"
	"time"
)

func ScheduleTask(ctx context.Context) {
	time.Sleep(untilWholePoint())
	for {
		select {
			case <-ctx.Done():
				logger.Info("Schedule service is shutting down...")
				return
			default:
				sendTask()
				time.Sleep(60 * time.Minute)
		}
	}
}

func sendTask() {
	current_date := strconv.Itoa(time.Now().Day())
	current_time := strconv.Itoa(time.Now().Hour())
	index := fmt.Sprintf("%s|%s", current_date, current_time)
	scan, collect, file := processSchedule()
	// StartScan
	if _, ok := scan[current_time]; ok {
		logger.Info("[SCHEDULED] " + time.Now().Format("2006-01-02 - 15:04:05") + " send task \"StartScan\"")
		go request.SendMissionToApi("StartScan", scan[current_time])
	}
	// StartCollect
	if _, ok := collect[index]; ok {
		logger.Info("[SCHEDULED] " + time.Now().Format("2006-01-02 - 15:04:05") + " send task \"StartCollect\"")
		go request.SendMissionToApi("StartCollect", collect[index])
	}
	// StartGetDrive
	if _, ok := file[index]; ok {
		logger.Info("[SCHEDULED] " + time.Now().Format("2006-01-02 - 15:04:05") + " send task \"StartGetDrive\"")
		go request.SendMissionToApi("StartGetDrive", file[index])
	}
}

func untilWholePoint() time.Duration {
	now := time.Now()
	next := now.Truncate(time.Hour).Add(time.Hour)
	return next.Sub(now)
}