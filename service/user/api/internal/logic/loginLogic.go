package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"api/internal/svc"
	"api/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func passwordMd5(password []byte) string {
	h := md5.New()
	h.Write([]byte(password))
	h.Write(secret)
	return hex.EncodeToString(h.Sum(nil))
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 实现登录功能
	// 1. 处理用户发来的请求 拿到用户名和密码
	// 2. 检验用户名密码 与数据库中的用户名密码 是否一致
	// 两种方式
	// 2.1 用 用户输入的用户名和密码(加密后) 去数据库中查询
	// select * from user where username=req.Username and password = req.Password
	// 2.2 用 用户名到数据库中查到结果 再判断密码
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err == sqlx.ErrNotFound {
		return &types.LoginResponse{
			Message: "用户名不存咋",
		}, nil
	}
	if err != nil {
		logx.Errorw("UserModel.FindOneByUsername", logx.Field("err", err))
		return &types.LoginResponse{
			Message: "用户名不存在",
		}, errors.New("内部错误")
	}
	if user.Password != passwordMd5([]byte(req.Password)) {
		return &types.LoginResponse{
			Message: "用户名或密码错误",
		}, nil
	}
	// 生成JWT
	now := time.Now().Unix()
	expire := l.svcCtx.Config.Auth.AccessExpire
	token, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, expire, user.UserId)
	if err != nil {
		logx.Errorw("l.getJwtToken faield", logx.Field("err", err))
		return nil, errors.New("内部错误")
	}
	// 3. 如果一致 登录成功 否则 登录失败
	return &types.LoginResponse{
		Message:      "登录成功",
		AccessToken:  token,
		AccessExpire: int(now + expire),
		RefreshAfter: int(now + (expire / 2)),
	}, nil
}

// 生成JWT的方法
func (l *LoginLogic) getJwtToken(secret string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	claims["author"] = "october" // 添加一些自定义的　key value 数据
	token := jwt.New(jwt.SigningMethodES256)
	token.Claims = claims
	return token.SignedString([]byte(secret))
}
