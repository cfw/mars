package httpx

import "reflect"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Ok() *Response {
	return new(Response)
}

func OkWithData(data interface{}) *Response {
	resp := Ok()
	if !reflect.ValueOf(data).IsNil() {
		resp.Data = data
	}
	return resp
}

func Error(c int, m string) *Response {
	return &Response{Code: c, Message: m}
}
