package taskservice

import (
	"edetector_API/pkg/mariadb/query"
	"strconv"
)

type task struct {
	taskId    string
	clientId  string
	status    int
}

func loadTasks(q task) []task {
	result := query.Load_stored_task(q.taskId, q.clientId, q.status)
	var tasks []task
	for _, v := range result {
		tmp := task{}
		tmp.taskId = v[0]
		tmp.clientId = v[1]
		tmp.status, _ = strconv.Atoi(v[2])
		tasks = append(tasks, tmp)
	}
	return tasks
}

func loadUnhandleTasks() []task {
	var q task
	q.taskId = "nil"
	q.clientId = "nil"
	q.status = 0
	return loadTasks(q)
}