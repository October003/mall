package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// CtxKey 自定义一个类型
type CtxKey string

const (
	// 使用自定义类型 声明context中存储的key，防止被他人覆盖
	CtxKeyAdminID CtxKey = "adminID"
)

// unaryInterceptor 客户端一元拦截器
func UnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println("客户端拦截器 in")
	// RPC 调用前
	// 编写客户端拦截器的逻辑
	adminID := ctx.Value(CtxKeyAdminID).(string)
	md := metadata.Pairs(
		"key1", "val1",
		"key1", "val1-2",
		"requestID", "123456",
		"token", "mall-order-october",
		"adminID", adminID, // 从外部获取，借助ctx上下文
	)
	ctx = metadata.NewOutgoingContext(ctx, md)
	err := invoker(ctx, method, req, reply, cc, opts...) // 实际的RPC调用
	// RPC 调用后
	fmt.Println("客户端拦截器 out")
	return err
}
