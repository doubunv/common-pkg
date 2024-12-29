package interceptors

import (
	"context"
	"github.com/doubunv/common-pkg/result/xcode"
	"github.com/zeromicro/go-zero/core/logc"
	"net/http"

	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ClientErrorInterceptor(appName, business string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		md.Set("app_name", appName)
		trace.Inject(ctx, otel.GetTextMapPropagator(), &md)
		ctx = metadata.NewOutgoingContext(ctx, md)

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			logc.Info(ctx, appName, method, req, reply)
			return nil
		}
		logc.Error(ctx, appName, method, req, reply, err.Error())

		gErr, ok := status.FromError(err)
		if !ok {
			return xcode.New(http.StatusInternalServerError, err.Error())
		}

		statusCode := gErr.Code()

		return xcode.New(int(statusCode), gErr.Message())
	}
}

func RetryDialOption() grpc.DialOption {
	retryPolicy := `{
		"methodConfig": [{
		  "retryPolicy": {
			  "MaxAttempts": 3,
			  "InitialBackoff": "1s",
			  "MaxBackoff": "1s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`

	return grpc.WithDefaultServiceConfig(retryPolicy)
}
