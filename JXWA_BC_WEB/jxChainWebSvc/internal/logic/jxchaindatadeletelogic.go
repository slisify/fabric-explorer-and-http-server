package logic

import (
	"context"

	"jxChainWebSvc/internal/svc"
	"jxChainWebSvc/internal/types"
	"jxChainWebSvc/fabricSdkOperate/fabricSdkApi"

	"github.com/tal-tech/go-zero/core/logx"
)

type JxChainDataDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJxChainDataDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) JxChainDataDeleteLogic {
	return JxChainDataDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JxChainDataDeleteLogic) JxChainDataDelete(req types.DeleteData) (*types.OnChainStatus, error) {
	var queryArgs = [][]byte{[]byte(req.AlertorID)}
	result := fabricSdkApi.DeleteData(queryArgs)
	if result == "ok" {
		return &types.OnChainStatus{
			Message: "ok",
		}, nil
	} else {
		return &types.OnChainStatus{
			Message: "failed",
		}, nil
	}
}
