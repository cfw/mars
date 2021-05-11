package httpx

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Ok() *Response {
	return new(Response)
}

func Error(c int, m string) *Response {
	return &Response{Code: c, Message: m}
}
