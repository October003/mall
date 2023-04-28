package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	// 数据库配置 除mysql外 还有MongoDB PostgreDB
	Mysql struct {
		// mysql 连接地址
		DataSource string
	}
	CacheRedis cache.CacheConf
	// jwt 鉴权配置
	Auth struct {
		AccessSecret string // jwt 密钥
		AccessExpire int64  // 有效期 单位: 秒
	}
}
