package interceptors

import (
	"context"
	"encoding/json"
	"github.com/doubunv/common-pkg/result/xcode"
	"github.com/zeromicro/go-zero/core/logc"
	"net/http"

	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type msgPrint struct {
	RpcName string      `json:"rpc_name"`
	Req     interface{} `json:"req"`
	Reply   interface{} `json:"reply"`
	Err     string      `json:"err"`
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
		}
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			marshal, _ := json.Marshal(msg)
			logc.Info(ctx, string(marshal))
			return nil
		}
		msg.Err = err.Error()
		marshal, _ := json.Marshal(msg)
		logc.Error(ctx, string(marshal))
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
