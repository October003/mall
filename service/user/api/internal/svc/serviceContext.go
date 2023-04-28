package svc

import (
	"api/internal/config"
	"api/internal/middleware"
	"mall/service/user/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config    config.Config
	Cost      rest.Middleware // 自定义路由中间件 字段名要与 api文件中的声明一致
	UserModel model.UserModel // 加入 user表 增删改查操作的Model
}

func NewServiceContext(c config.Config) *ServiceContext {
	// UserModel 接口类型
	// defaultUserModel 实现了该接口
	// 调用构造函数得到  *defaultUserModel
	// NewUserModel(sqlx.conn)
	// 需要sqlx.conn 的数据库连接
	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(sqlxConn, c.CacheRedis),
		Cost:      middleware.NewCostMiddleware().Handle,
	}
}
