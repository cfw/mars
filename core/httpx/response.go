package httpx

import (
	"github.com/gogo/protobuf/proto"
	"reflect"
)

type Response struct {
	Code    int           `json:"code"`
	Message string        `json:"message,omitempty"`
	Data    proto.Message `json:"data,omitempty"`
}

func Ok() *Response {
	return new(Response)
}

func OkWithData(data proto.Message) *Response {
	resp := Ok()
	if d, ok := data.(interface {
		Size() int
	}); ok {
		if d.Size() > 0 {
			resp.Data = data
		}
	} else {
		if !reflect.ValueOf(data).IsNil() {
			resp.Data = data
		}
	}
	return resp
}

func Error(c int, m string) *Response {
	return &Response{Code: c, Message: m}
}
