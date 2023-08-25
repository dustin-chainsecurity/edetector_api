# API Service Documentation
This service acts as a bridge that enables seamless communication, data transfer and interaction between the user-facing frontend and the robust backend server.

### Run Service
```console
go run cmd/api/api.go 5000
```

### Monitor Logs
You can view service logs from the console, `./cmd/api/*.log`, and `/var/log/syslog`

# Version
### v0.0.1 (2023.08.21)
Functioning APIs :
- Backend-Related : SaveAgent, AddDevice, UpdateProgress, Test
- Member-Related : Login, LoginWithToken, Signup
- Search Evidence Page
- Analysis Page
- Group Settings

### v1.0.0 (2023.08.24)
Enhancements :
- Added settings-related APIs
- Integrated zap logger with syslog

---

# REST API

## Member

<details>
<summary> <code> <b>POST</b> /login </code> </summary>
<br/>

Request
```json
"Body": {
    "username": "example",
    "password": "example"
}
```
Response
```json
"Body": {
    "success": true,
    "message": "success",
    "user": {
        "username": "example",
        "token": "token"
    }
}
```
</details>

<details>
<summary> <code> <b>POST</b> /loginWithToken </code> </summary>
<br/>

Request
```json
"Body": {
    "token": "token"
}
```
Response
```json
"Body": {
    "success": true,
    "message": "success",
    "user": {
        "username": "example",
        "token": "token"
    }
}
```
</details>

## Search Evidence Page

<details>
<summary> <code> <b>GET</b> /searchEvidence/detectDevices </code> </summary>
<br/>

Request
```json
"Header": {"Authorization": "token"}
```
Response
```json
{
    "isSuccess": true,
    "data": [
        {
            "deviceId": "example",
            "connection": true,
            "innerIP": "example",
            "deviceName": "example",
            "groups": ["example"],
            "detectionMode": true,
            "scanSchedule": ["example"],
            "scanFinishTime": {
                "isFinish": true,
                "progress": 0,
                "finishTime": 0
            },
            "collectSchedule": {
                "date": "example",
                "time": "example"
            },
            "collectFinishTime": {
                "isFinish": true,
                "progress": 0,
                "finishTime": 0
            },
            "fileDownloadDate": {
                "date": "example",
                "time": "example"
            },
            "fileFinishTime": {
                "isFinish": true,
                "progress": 0,
                "finishTime": 0
            },
            "imageFinishTime": {
                "isFinish": true,
                "progress": 0,
                "finishTime": 0
            }
        }
    ]
}
```
</details>

<details>
<summary> <code> <b>POST</b> /searchEvidence/refresh </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>POST</b> /task/sendMission </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>POST</b> /task/detectionMode </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

## Analysis Page

<details>
<summary> <code> <b>GET</b> /analysisPage/allDeviceDetail </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>POST</b> /analysisPage/template </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>GET</b> /analysisPage/template </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>GET</b> /analysisPage/template/:id </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>PUT</b> /analysisPage/template/:id </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>DELETE</b> /analysisPage/template/:id </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

## Setting Page
### Group Settings

<details>
<summary> <code> <b>POST</b> /group </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>GET</b> /group </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>GET</b> /group/:id </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>PUT</b> /group/:id </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>DELETE</b> /group/:id </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>POST</b> /group/device </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>DELETE</b> /group/device </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

### Other Settings

<details>
<summary> <code> <b>GET</b> /setting/:field </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>

<details>
<summary> <code> <b>POST</b> /setting/:field </code> </summary>
<br/>

Request
```json

```
Response
```json

```
</details>
