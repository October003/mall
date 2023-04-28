package logic

import (
	"context"
	"errors"

	"mall/service/user/rpc/internal/svc"
	"mall/service/user/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserReq) (*user.GetUserResp, error) {
	// todo: add your logic here and delete this line
	// 根据userID 查询数据库返回用户信息
	one, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, in.UserID)
	// 1. 查数据失败
	if errors.Is(err, sqlx.ErrNotFound) {
		return nil, errors.New("无效的userID")
	}
	// 2. 根据 userID 查不到用户记录
	if err != nil {
		logx.Errorw("user.rpc.GetUser FindOneByuserId failed", logx.Field("err", err))
		return nil, errors.New("查询失败")
	}
	return &user.GetUserResp{
		Username: one.Username,
		Gender:   one.Gender,
	}, nil
}
