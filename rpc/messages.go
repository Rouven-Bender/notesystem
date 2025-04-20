package rpc

type BaseMessage struct {
	Method string `json:"method"`
}

type EasterMessage BaseMessage

type BaseResponse struct {
	Response string `json:"response"`
}

type BaseError struct {
	Error string `json:"error"`
}
