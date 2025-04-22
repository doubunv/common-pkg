package interceptors

import (
	"context"
	"github.com/doubunv/common-pkg/result/xcode"
	"github.com/zeromicro/go-zero/core/logc"
	"google.golang.org/grpc/status"
	"net/http"

	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type msgPrint struct {
	RpcName string      `json:"rpc_name"`
	Req     interface{} `json:"req"`
	Reply   interface{} `json:"reply"`
	Err     string      `json:"err"`
	Method  string      `json:"method"`
}

func ClientInterceptor(rpcName string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		md.Set("rpc_name", rpcName)
		trace.Inject(ctx, otel.GetTextMapPropagator(), &md)
		ctx = metadata.NewOutgoingContext(ctx, md)
		msg := msgPrint{
			RpcName: rpcName,
			Req:     req,
			Reply:   reply,
			Method:  method,
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			logc.Infof(ctx, "%+v", msg)
			return nil
		}
		msg.Err = err.Error()
		logc.Errorf(ctx, "%+v", msg)
		gErr, ok := status.FromError(err)
		if !ok {
			return xcode.New(http.StatusInternalServerError, "Service catch err")
		}

		statusCode := http.StatusInternalServerError
		return xcode.New(statusCode, gErr.Message())
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
