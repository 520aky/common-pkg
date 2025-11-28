package rpcclient

import (
	"fmt"
	"github.com/zeromicro/go-zero/zrpc"
)

func GenRpcTarget(hosts string) string {
	return fmt.Sprintf("%s", hosts)
}

type Config struct {
	Host    string
	RpcName string
	Timeout int64
}

func MustNewClient(conf Config) zrpc.Client {
	return zrpc.MustNewClient(
		zrpc.RpcClientConf{
			Timeout: conf.Timeout, //10s
			Target:  GenRpcTarget(conf.Host),
		},
		zrpc.WithUnaryClientInterceptor(ClientInterceptor(conf.RpcName)),
		zrpc.WithDialOption(RetryDialOption()),
	)
}

func NewClient(conf Config) (zrpc.Client, error) {
	return zrpc.NewClient(
		zrpc.RpcClientConf{
			Timeout: conf.Timeout, //10s
			Target:  GenRpcTarget(conf.Host),
		},
		zrpc.WithUnaryClientInterceptor(ClientInterceptor(conf.RpcName)),
		zrpc.WithDialOption(RetryDialOption()),
	)
}
