# Task Service Documentation
This service schedules and allocates tasks for the working server in order to smoothen the overall working process.

### Run Service
```console
go run cmd/taskservice/taskservice.go
```
### Monitor Logs
You can view service logs from the console, `./cmd/taskservice/*.log`, and `/var/log/syslog`

# Version
### v1.0.0 (2023.08.21)
Functions :
- Task allocation for working server
- Task scheduling for working server

### v1.0.1 (2023.08.24)
Enhancements :
- Integrated zap logger with syslog

---