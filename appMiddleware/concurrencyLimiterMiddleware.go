package appMiddleware

import (
	"github.com/zeromicro/go-zero/core/syncx"
	"net/http"
)

type ConcurrencyLimiter struct {
	limiter syncx.Limit
}

func NewConcurrencyLimiter(n int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		limiter: syncx.NewLimit(n),
	}
}

func (m *ConcurrencyLimiter) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.TryBorrow() {
			// 直接拒绝
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"code":429,"msg":"too many requests"}`))
			return
		}
		defer m.limiter.Return()
		next(w, r)
	}
}
