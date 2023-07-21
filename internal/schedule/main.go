package schedule

import (
	"context"
	"edetector_API/internal/request"
	"fmt"
	"strconv"
	"time"
)

func ScheduleTask(ctx context.Context) {
	time.Sleep(untilWholePoint())
	for {
		select {
			case <-ctx.Done():
				fmt.Println("Schedule service is shutting down...")
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
		fmt.Println("send task \"StartScan\"")
		go request.SendMissionToApi("StartScan", scan[current_time])
	}
	// StartCollect
	if _, ok := collect[index]; ok {
		fmt.Println("send task \"StartCollect\"")
		go request.SendMissionToApi("StartCollect", collect[index])
	}
	// StartGetDrive
	if _, ok := file[index]; ok {
		fmt.Println("send task \"StartGetDrive\"")
		go request.SendMissionToApi("StartGetDrive", file[index])
	}
}

func untilWholePoint() time.Duration {
	now := time.Now()
	next := now.Truncate(time.Hour).Add(time.Hour)
	return next.Sub(now)
}