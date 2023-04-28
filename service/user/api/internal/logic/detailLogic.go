package logic

import (
	"context"
	"errors"
	"fmt"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailRequest) (resp *types.DetailResponse, err error) {
	// jwt 鉴权后  解析出来的数据
	fmt.Printf("JWT userId:%v\n",l.ctx.Value("userId"))
	fmt.Printf("JWT author:%v\n",l.ctx.Value("author"))
	// 1. 拿到用户的UesrID
	if req.UserID < 0 {
		return nil, errors.New("参数错误")
	}
	// 2. 通过UseID 到数据库中查询 用户的数据
	user, err := l.svcCtx.UserModel.FindOneByUserId(l.ctx, req.UserID)
	// 查数据失败
	if err != nil {
		if err != sqlx.ErrNotFound {
			logx.Errorw("UserModel.FindOneByUserId failed", logx.Field("err", err))
			return nil, errors.New("内部错误")
		}
		return nil, errors.New("用户不存在")
	}
	// 3. 格式化数据(数据库里存的字段和前端要求的字段不太一致) 将用户的具体信息 返回
	return &types.DetailResponse{
		Username: user.Username,
		Gender:   int(user.Gender),
		Message:  "查询用户成功",
	}, nil
}
