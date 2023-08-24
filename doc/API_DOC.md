# API Service Documentation
This service acts as a bridge that enables seamless communication, data transfer and interaction between the user-facing frontend and the robust backend server.

### Run Service
```console
go run cmd/api/api.go 5000
```

### Monitor Logs
You can view service logs from the console, `./cmd/api/*.log`, and `/var/log/syslog`

# REST API

## Member

<details>
<summary> <code> <b>POST</b> /login </code> </summary>
<br/>

Request
```json
Body: {
    "username": string,
    "password": string
}
```
Response
```json
Body: {
    "success": boolean,
    "message": string,
    "user": {
        "username": string,
        "token": string
    }
}
```
</details>

<details>
<summary> <code> <b>POST</b> /loginWithToken </code> </summary>
<br/>

Request
```json
Body: {
    "token": string
}
```
Response
```json
Body: {
    "success": boolean,
    "message": string,
    "user": {
        "username": string,
        "token": string
    }
}
```
</details>

## Search Evidence Page

<details>
<summary> <code> <b>GET</b> /searchEvidence/detectDevices </code> </summary>
<br/>

Request
```console
Header: {"Authorization": "token"}
```
Response
```json
{
    "isSuccess": boolean,
    "data": [
        {
            "deviceId": string,
            "connection": boolean,
            "innerIP": string,
            "deviceName": string,
            "groups": [],
            "detectionMode": boolean,
            "scanSchedule": []string,
            "scanFinishTime": {
                "isFinish": boolean,
                "progress": int,
                "finishTime": int
            },
            "collectSchedule": {
                "date": string,
                "time": string
            },
            "collectFinishTime": {
                "isFinish": boolean,
                "progress": int,
                "finishTime": int
            },
            "fileDownloadDate": {
                "date": string,
                "time": string
            },
            "fileFinishTime": {
                "isFinish": boolean,
                "progress": int,
                "finishTime": int
            },
            "imageFinishTime": {
                "isFinish": boolean,
                "progress": int,
                "finishTime": int
            }
        }, ...
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
