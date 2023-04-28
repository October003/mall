package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	UserRPC zrpc.RpcClientConf // 连接其它微服务的rpc客户端
}
