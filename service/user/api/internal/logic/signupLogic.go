package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"mall/service/user/model"
	"time"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var secret = []byte("october")

type SignupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignupLogic {
	return &SignupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SignupLogic) Signup(req *types.SignupRequest) (resp *types.SignupResponse, err error) {
	// 参数校验
	if req.RePassword != req.Password {
		return nil, errors.New("两次输入的密码不一致")
	}
	logx.Debugv(req) // json.Marshal(req)
	logx.Debugf("req:%#v\n", req)
	// 0.查询uesrname 是否已经被注册
	u, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	// 0.1 查询数据库失败了
	// var ErrNotFound = sql.ErrNoRows
	if err != nil && err != sqlx.ErrNotFound {
		logx.Errorw("user_signup_UserModel.FindOneByUsername failed",
			logx.Field("err", err),
		)
		fmt.Printf("FindByUsername err:%v\n", err)
		return nil, errors.New("内部错误")
	}
	// 0.2 查到记录  表示该用户已经被注册
	if u != nil {
		return nil, errors.New("用户名已存在")
	}

	// 2.加密密码(加盐 MD5)
	h := md5.New()
	h.Write([]byte(req.Password))
	h.Write(secret)
	passwordStr := hex.EncodeToString(h.Sum(nil))
	// 在这里写你的业务逻辑
	fmt.Printf("req:%#v\n", req)
	// 把用户的注册信息 保存到数据库中
	// 1.生成UserId (雪花算法)
	// 2.加密密码(加盐 MD5)
	user := &model.User{
		UserId:   time.Now().Unix(), //这里简化
		Username: req.Username,
		Password: passwordStr, //不能存明文
		Gender:   1,
	}
	if _, err := l.svcCtx.UserModel.Insert(context.Background(), user); err != nil {
		logx.Errorf("user_signup_UserModel.Insert failed,err:%v\n,err")
		return nil, err
	}
	return &types.SignupResponse{Message: "success"}, nil
}
