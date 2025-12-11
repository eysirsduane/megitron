// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transfer

import (
	"context"

	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferUsdtCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferUsdtCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferUsdtCreateLogic {
	return &TransferUsdtCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferUsdtCreateLogic) TransferUsdtCreate(req *types.TransferUsdtCreateReq) (resp *types.TransferUsdtCreateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
