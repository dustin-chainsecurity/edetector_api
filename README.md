# eDetector API Documentation
This repository contains three microservices : 

- **Web API Service**<br />
This service acts as a bridge that enables seamless communication, data transfer and interaction between the user-facing frontend and the robust backend server. <br />See [] for more details.

- **Websocket Service**<br />
This service is designed to enable the working server to inform the frontend to refresh device details if the worker tasks' status has changed, providing real-time user experience. <br />See [] for more details.

- **Task Service**<br />
This service schedules and allocates tasks for the working server in order to smoothen the overall working process. <br />See [] for more details.

## Directory Layout

Brief description of the layout :

* `README.md` contains detailed description of the repository.
* `docs`contains detailed documentations for the services.
* `config` contains environment variables to be used in the repository.
* `cmd` contains main packages (entrypoints), each subdirectory of `cmd` is a main package.
* `pkg` holds common packages that can be shared among different repositories.
* `internal` holds internal packages only for this repository.
* `api` contains packages for each api route.
* `test` holds testing programs.

## Getting Started
### Requirements
This repository requires the following environment : 
- Operating System : Ubuntu 22.04.2 LTS
- Go Version : 1.18.1 linux/amd64 (or above but you have to change go.mod)
    ```bash
    # version check
    go version
    ```
- rsyslog installed and enabling tcp input

### Installation
Run the commands below to get all environments set : 
```bash
git clone https://github.com/yu-niverse/edetector_api.git
cd edetector_api
go mod download
```

---

