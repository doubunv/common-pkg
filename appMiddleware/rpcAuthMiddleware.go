package appMiddleware

import (
	"context"
	"fmt"
	"gitlab.888bbm.com/go-package/common-pkg/result/xcode"
	"net/http"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc"
)

type RpcAuthMiddleware struct {
	Debug bool
}

func NewRpcAuthMiddleware() *RpcAuthMiddleware {
	return &RpcAuthMiddleware{}
}

func (m *RpcAuthMiddleware) Handle() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if p := recover(); p != nil {
				logc.Error(ctx, err, string(debug.Stack()))
				err = xcode.NewRpc(http.StatusInternalServerError, "Rpc Server error.")
				return
			}
		}()

		logc.Info(ctx, info.FullMethod+",RpcRequest:", req)

		//mdData, ok := metadata.FromIncomingContext(ctx)
		//if !ok {
		//	return nil, errors.New("no rpc auth metadata. ")
		//}
		//ctx = metadata.NewOutgoingContext(ctx, mdData)
		//ctx = rpc.ContextMetadataInLog(ctx)

		resp, err = handler(ctx, req)
		if err != nil {
			logc.Error(ctx, fmt.Sprintf(info.FullMethod+",rpc错误信息：%v", err))
		} else {
			logc.Info(ctx, info.FullMethod+",RpcResponse:", resp)
		}
		return resp, err
	}
}
