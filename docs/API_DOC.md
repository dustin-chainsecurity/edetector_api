# API Service Documentation
This service acts as a bridge that enables seamless communication, data transfer and interaction between the user-facing frontend and the robust backend server.

### Run Service
```console
go run cmd/api/api.go 5000
```

### Monitor Logs
You can view service logs from the console, `./cmd/api/*.log`, and `/var/log/syslog`

# Version
### v1.0.0
Functioning APIs :
- Backend-Related : SaveAgent, AddDevice, UpdateProgress, Test
- Member-Related : Login, LoginWithToken, Signup
- Search Evidence Page
- Analysis Page
- Group Settings

### v1.0.1
Enhancements :
- Added settings keyword API
- Integrated zap logger with syslog

### v1.0.2
Enhancements :
- Integrated gin logger with zap logger
- Able to handle more task status (terminated and failure cases)

Fixes :
- Changed API return status and body messages when error occurs
- Added null string check for group names and template names
- Checked terminating status before adding tasks

### v1.0.3 
Enhancements : 
- Added settings API (logs, whitelist, blacklist, hacklist, user, keyImage)
- Logging to database
- Added version information for devices

---
