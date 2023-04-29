package svc

import (
	"mall/service/order/api/internal/config"
	"mall/service/order/api/internal/interceptor"
	"mall/service/order/model"
	"mall/service/user/rpc/userclient" // RPC 客户端代码

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRPC    userclient.User // RPC 客户端
	OrderModel model.OrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		UserRPC: userclient.NewUser(
			zrpc.MustNewClient(
				c.UserRPC,
				zrpc.WithUnaryClientInterceptor(interceptor.UnaryInterceptor),
			),
		),
	}
}
