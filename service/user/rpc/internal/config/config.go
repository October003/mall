package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// 数据库设置
	Mysql struct {
		// 数据库连接地址
		DataSource string
	}
	// redis
	CacheRedis cache.CacheConf
}
