package logic

import (
	"context"
	"encoding/json"

	"jxChainWebSvc/fabricSdkOperate/fabricSdkApi"
	"jxChainWebSvc/internal/svc"
	"jxChainWebSvc/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type JxChainDataGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJxChainDataGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) JxChainDataGetLogic {
	return JxChainDataGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// func (l *JxChainDataGetLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
// 	claims := make(jwt.MapClaims)
// 	claims["exp"] = iat + seconds
// 	claims["iat"] = iat
// 	claims["userId"] = userId
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	token.Claims = claims
// 	return token.SignedString([]byte(secretKey))
// }

func (l *JxChainDataGetLogic) JxChainDataGet(req types.QueryData) (*types.OnChainDataPut, error) {
	//调用链码进行查询
	var alertData types.OnChainDataPut
	var queryArgs = [][]byte{[]byte(req.AlertorID)}
	alertDataJson := fabricSdkApi.QueryData(queryArgs)
	err := json.Unmarshal(alertDataJson, &alertData)
	if err != nil {
		panic("Unmarshal failed!")
	}
	//进行鉴权
	// // ---start---
	// now := time.Now().Unix()
	// accessExpire := l.svcCtx.Config.Auth.AccessExpire
	// jwtToken, err := l.getJwtToken(l.svcCtx.Config.Auth.AccessSecret, now, l.svcCtx.Config.Auth.AccessExpire, userInfo.Id)
	// if err != nil {
	// 	return nil, err
	// }
	// // ---end---

	return &alertData, nil
}
