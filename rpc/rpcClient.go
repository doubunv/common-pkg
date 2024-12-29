package rpc

import (
	"fmt"
	"github.com/doubunv/common-pkg/rpc/interceptors"

	"github.com/zeromicro/go-zero/zrpc"
)

func GenRpcTarget(hosts string) string {
	return fmt.Sprintf("%s", hosts)
}

type Config struct {
	Host    string
	RpcName string
}

func MustNewClient(conf Config) zrpc.Client {
	return zrpc.MustNewClient(
		zrpc.RpcClientConf{
			Target: GenRpcTarget(conf.Host),
		},
		zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor(conf.RpcName)),
		zrpc.WithDialOption(interceptors.RetryDialOption()),
	)
}

func NewClient(conf Config) (zrpc.Client, error) {
	return zrpc.NewClient(
		zrpc.RpcClientConf{
			Target: GenRpcTarget(conf.Host),
		},
		zrpc.WithUnaryClientInterceptor(interceptors.ClientErrorInterceptor(conf.RpcName)),
		zrpc.WithDialOption(interceptors.RetryDialOption()),
	)
}
