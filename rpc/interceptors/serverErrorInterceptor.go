package interceptors

import (
	"context"
	"gitlab.888bbm.com/go-package/common-pkg/logs/xcode"

	"google.golang.org/grpc"
)

func ServerErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, xcode.FromError(err).Err()
	}
}
