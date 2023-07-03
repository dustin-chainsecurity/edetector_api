package api

type Response struct {
	IsSuccess  bool           `json:"isSuccess"`
	Data       interface{}    `json:"Data"`
}