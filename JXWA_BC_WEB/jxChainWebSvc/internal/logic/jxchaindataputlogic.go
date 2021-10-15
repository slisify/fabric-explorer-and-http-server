package logic

import (
	"context"

	"jxChainWebSvc/fabricSdkOperate/fabricSdkApi"
	"jxChainWebSvc/internal/svc"
	"jxChainWebSvc/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type JxChainDataPutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJxChainDataPutLogic(ctx context.Context, svcCtx *svc.ServiceContext) JxChainDataPutLogic {
	return JxChainDataPutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JxChainDataPutLogic) JxChainDataPut(req types.OnChainDataPut) (*types.OnChainStatus, error) {
	// todo: add your logic here and delete this line
	//调用SDK链码功能，进行数据的重新组装，然后调用对应的API即可
	// fmt.Println("alertor id is %s", req.AlertorID)
	// fmt.Println("alertor id is %s", req.AlertTime)
	// fmt.Println("alertor id is %s", req.HouseNumber)
	// fmt.Println("alertor id is %s", req.HouseOwnerName)
	var queryArgs = [][]byte{[]byte(req.AlertorID), []byte(req.AlertTime), []byte(req.HouseNumber), []byte(req.HouseOwnerName)}
	result, err := fabricSdkApi.ExecuteChainCode(queryArgs)
	if err == nil {
		return &types.OnChainStatus{
			Message: result,
		}, nil
	} else {
		return &types.OnChainStatus{
			Message: "put data failed",
		}, nil
	}
	return &types.OnChainStatus{
		Message: "put  ok!!!!!",
	}, nil
}
