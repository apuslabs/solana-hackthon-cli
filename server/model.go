package server

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Health struct {
	Busy bool   `json:"busy"`
	Port string `json:"port"`
}
