package appMiddleware

import (
	"github.com/doubunv/common-pkg/result"
	"net/http"
)

type MiddlewareOption func(m *UserAgentMiddleware)

func WithCheckOption(check CheckRequestTokenFunc) MiddlewareOption {
	return func(m *UserAgentMiddleware) {
		m.check = check
	}
}

type UserAgentMiddleware struct {
	check        CheckRequestTokenFunc
	noVerifyPath map[string]int
}

func NewUserAgentMiddleware(noVerifyPath map[string]int, ops ...MiddlewareOption) *UserAgentMiddleware {
	res := &UserAgentMiddleware{
		noVerifyPath: noVerifyPath,
	}
	for _, op := range ops {
		op(res)
	}
	return res
}

func (m *UserAgentMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newReq, err := MustAuthTokenRequest(r, m.check, m.noVerifyPath)
		if err != nil {
			result.HttpErrorResult(r.Context(), w, err)
			return
		}
		next(w, newReq)
	}
}
