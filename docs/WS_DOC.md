# Websocket Service Documentation
This service is designed to enable the working server to inform the frontend to refresh device details if the worker tasks' status has changed, providing real-time user experience.

### Run Service
```console
go run cmd/websocket/websocket.go 5050
```

### Monitor Logs
You can view service logs from the console, `./cmd/websocket/*.log`, and `/var/log/syslog`

# Version
### v1.0.0 (2023.08.21)
Functions :
- Informs frontend to refresh page when task status have been updated
- Scheduling task APIs
- Update Task API for backend usage

### v1.0.1 (2023.08.24)
Enhancements :
- Integrated zap logger with syslog

### v1.0.2 (2023.09.01)
Enhancements :
- Integrated gin logger with zap logger
Fixes:
- Changed API return status and body messages when error occurs

---

## REST API
