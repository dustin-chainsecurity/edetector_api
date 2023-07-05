package task

type TaskRequest struct {
	Key      string      `json:"key"`
	User     string      `json:"user"`
	Message  interface{} `json:"message"`
}

type TaskPacket struct {
	Key      string      `json:"key"`
	Work     string      `json:"work"`
	User     string      `json:"user"`
	Message  interface{} `json:"message"`
}

type TaskResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}