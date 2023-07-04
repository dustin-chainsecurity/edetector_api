package task

type TaskPacket struct {
	Key      string      `json:"Key"`
	Work     string      `json:"Work"`
	User     string      `json:"User"`
	Message  interface{} `json:"Message"`
}

type TaskResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}