package result

import (
	"net/http"
)

type ResponseBean struct {
	Code  int         `json:"code"`
	Msg   string      `json:"message"`
	Data  interface{} `json:"data"`
	Trace string      `json:"trace"`
}

func Success(data interface{}, trace string) *ResponseBean {
	if data == nil {
		data = struct{}{}
	}
	return &ResponseBean{http.StatusOK, "success", data, trace}
}

func Error(errCode int, errMsg string, trace string) *ResponseBean {
	return &ResponseBean{errCode, errMsg, struct{}{}, trace}
}
