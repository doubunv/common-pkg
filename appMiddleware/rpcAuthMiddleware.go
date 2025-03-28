package appMiddleware

import (
	"context"
	"errors"
	"github.com/doubunv/common-pkg/headInfo"
	"github.com/doubunv/common-pkg/result/xcode"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
	"net/http"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc"
)

type RpcAuthMiddleware struct {
	Debug bool
}

type PathWhiteHandle func(ctx context.Context) map[string]int

func NewRpcAuthMiddleware() *RpcAuthMiddleware {
	return &RpcAuthMiddleware{}
}

func (m *RpcAuthMiddleware) contextMetadataInLog(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	var fieldList = make([]logx.LogField, 0)
	for k, v := range md {
		fieldList = append(fieldList, logx.Field(k, v))
	}
	ctxNew := logx.ContextWithFields(ctx, fieldList...)

	return ctxNew
}

func (m *RpcAuthMiddleware) Handle(getPathWhite PathWhiteHandle) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if p := recover(); p != nil {
				logc.Error(ctx, err, string(debug.Stack()))
				err = xcode.NewRpc(http.StatusInternalServerError, "Rpc Server error.")
				return
			}
		}()

		mdData, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = metadata.NewOutgoingContext(ctx, mdData)
			ctx = m.contextMetadataInLog(ctx)
			if headInfo.GetBusinessCode(ctx) == "" {
				if getPathWhite == nil {
					return nil, errors.New("Rpc header data err")
				}
				whitePath := getPathWhite(ctx)
				if _, okPath := whitePath[info.FullMethod]; !okPath {
					return nil, errors.New("Rpc header data err")
				}
			}
		}

		resp, err = handler(ctx, req)
		return resp, err
	}
}
