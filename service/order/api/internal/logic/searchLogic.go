package logic

import (
	"context"
	"errors"
	"strconv"

	"mall/service/order/api/internal/interceptor"
	"mall/service/order/api/internal/svc"
	"mall/service/order/api/internal/types"
	"mall/service/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	// todo: add your logic here and delete this line
	// 1.根据请求参数中的订单号查询数据库中的订单记录
	orderID, _ := strconv.Atoi(req.OrderID)
	one, err := l.svcCtx.OrderModel.FindOne(l.ctx, int64(orderID))
	if err != nil {
		return nil, errors.New("内部错误")
	}
	// 如何存入adminID？
	l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdminID, "33")
	// 2.根据订单记录中的 user_id 去查询用户数据 (通过RPC调用user服务)
	userResp, err := l.svcCtx.UserRPC.GetUser(l.ctx, &userclient.GetUserReq{UserID: one.UserId}) //1682424219
	if err != nil {
		logx.Errorw("UserRPC.GetUser failed", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}
	// 3.拼接返回结果(因为我们这个接口的数据不是由我们一个服务组成的)
	return &types.SearchResponse{
		OederID:  "1682424219",
		Status:   100,                    //根据实际查询的订单记录来复制，这里写的是假数据
		Username: userResp.GetUsername(), // RPC 调用user服务拿到的数据
	}, nil
}
