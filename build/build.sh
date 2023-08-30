#!/bin/bash
source config/app.env
go build -o build/api_$API_VERSION.exe -ldflags "-X main.version=$API_VERSION" cmd/api/api.go
go build -o build/taskservice_$WS_VERSION.exe -ldflags "-X main.version=$WS_VERSION" cmd/taskservice/taskservice.go
go build -o build/websocket_$TASK_VERSION.exe -ldflags "-X main.version=$TASK_VERSION" cmd/websocket/websocket.go

