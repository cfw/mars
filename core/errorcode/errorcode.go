package errorcode

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

const (
	codeKey = "code"
	MsgKey  = "message"
)

type ErrorCode struct {
	Code    int
	Message string
}

func Error(e ErrorCode) error {
	return GrpcStatusError(e.Code, e.Message)
}

func GrpcStatusError(code int, msg string) error {
	st := status.New(codes.InvalidArgument, "")
	br := &errdetails.ErrorInfo{Metadata: map[string]string{codeKey: strconv.Itoa(code), MsgKey: msg}}
	st, _ = st.WithDetails(br)
	return st.Err()
}

func Convert(err error) (int, string) {
	s := status.Convert(err)
	if s.Code() == codes.InvalidArgument {
		var c int
		var m string
		for _, detail := range s.Details() {
			switch t := detail.(type) {
			case *errdetails.ErrorInfo:
				data := t.GetMetadata()
				c, _ = strconv.Atoi(data[codeKey])
				m = data[MsgKey]
			}
		}
		return c, m
	}
	return int(s.Code()), s.Message()
}
