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
			Timeout:       10000, //10s
			Target:        GenRpcTarget(conf.Host),
			NonBlock:      true,
			KeepaliveTime: 30000,
		},
		zrpc.WithUnaryClientInterceptor(interceptors.ClientInterceptor(conf.RpcName)),
		zrpc.WithDialOption(interceptors.RetryDialOption()),
	)
}

func NewClient(conf Config) (zrpc.Client, error) {
	return zrpc.NewClient(
		zrpc.RpcClientConf{
			Timeout:       10000, //10s
			Target:        GenRpcTarget(conf.Host),
			NonBlock:      true,
			KeepaliveTime: 30000,
		},
		zrpc.WithUnaryClientInterceptor(interceptors.ClientInterceptor(conf.RpcName)),
		zrpc.WithDialOption(interceptors.RetryDialOption()),
	)
}
