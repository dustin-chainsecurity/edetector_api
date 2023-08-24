# Task Service Documentation
This service schedules and allocates tasks for the working server in order to smoothen the overall working process.

### Run Service
```console
go run cmd/taskservice/taskservice.go <forward port>
```
### Monitor Logs
You can view service logs from the console, `./cmd/taskservice/*.log`, and `/var/log/syslog`