package query

type Log struct {
	Id      int    `json:"id"`
	Level   string `json:"level"`
	Service string `json:"service"`
	Message string `json:"content"`
}


