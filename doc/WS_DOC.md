# Websocket Service Documentation
This service is designed to enable the working server to inform the frontend to refresh device details if the worker tasks' status has changed, providing real-time user experience.

### Run Service
```console
go run cmd/websocket/websocket.go 5050
```

### Monitor Logs
You can view service logs from the console, `./cmd/websocket/*.log`, and `/var/log/syslog`

## API List
