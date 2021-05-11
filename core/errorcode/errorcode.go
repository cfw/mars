package errorcode

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

type ErrorCode struct {
	Code    int
	Message string
}

func Error(e ErrorCode) error {
	return GrpcStatusError(e.Code, e.Message)
}

func GrpcStatusError(code int, msg string) error {
	st := status.New(codes.Unknown, "")
	br := &errdetails.ErrorInfo{Metadata: map[string]string{"code": strconv.Itoa(code), "message": msg}}
	st, _ = st.WithDetails(br)
	return st.Err()
}
func Convert(err error) (int, string) {
	s := status.Convert(err)
	var c int
	var m string
	for _, detail := range s.Details() {
		switch t := detail.(type) {
		case *errdetails.ErrorInfo:
			data := t.GetMetadata()
			c, _ = strconv.Atoi(data["code"])
			m = data["message"]
		}
	}
	return c, m
}
