// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transfer

import (
	"context"

	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferTrxCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferTrxCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferTrxCreateLogic {
	return &TransferTrxCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferTrxCreateLogic) TransferTrxCreate(req *types.TransferTrxCreateReq) (resp *types.TransferTrxCreateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
