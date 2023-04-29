package main

import (
	"context"
	"flag"
	"fmt"

	"mall/service/user/rpc/internal/config"
	"mall/service/user/rpc/internal/server"
	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	// 注册 server端 拦截器
	s.AddUnaryInterceptors()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start() // 启动RPC服务
}

func MyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 调用前
	fmt.Println("服务端拦截器 in")
	// 拦截器的业务逻辑
	// 取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "need metadata")
	}
	fmt.Printf("metadata:%#v\n", md)
	// 根据metadata中的数据 进行一些检验处理
	if md["token"][0] != "mall-order-october" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	m, err := handler(ctx, req) // 实际的RPC方法调用
	// 调用后
	fmt.Println("服务端拦截器 out")
	return m, err
}
