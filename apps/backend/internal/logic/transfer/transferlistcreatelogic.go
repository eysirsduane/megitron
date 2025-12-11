// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transfer

import (
	"context"

	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferListCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferListCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferListCreateLogic {
	return &TransferListCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferListCreateLogic) TransferListCreate(req *types.TransferListReq) (resp *types.TransferListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
