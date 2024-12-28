package xcode

import (
	"errors"
	"gitlab.888bbm.com/go-package/common-pkg/result/xcode"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrHandler(err error) (int, any) {
	var (
		code int
		msg  string
	)
	var xErr xcode.XCode
	if errors.As(err, &xErr) {
		code = xErr.Code()
		msg = xErr.Error()
	} else {
		code = http.StatusInternalServerError
	}

	type Response struct {
		Code int         `json:"code"`
		Msg  string      `json:"message"`
		Data interface{} `json:"data"`
	}
	return http.StatusOK, &Response{
		Code: code,
		Msg:  msg,
		Data: struct{}{},
	}
}

func FromError(err error) *status.Status {
	var code xcode.XCode
	if errors.As(err, &code) {
		return status.New(codes.Code(code.Code()), code.Error())
	}

	return status.New(codes.Code(http.StatusInternalServerError), err.Error())
}

// Deprecated : 使用 pkg/xcode/New 替换
func New(code int, msg string) xcode.Code {
	return xcode.New(code, msg)
}
