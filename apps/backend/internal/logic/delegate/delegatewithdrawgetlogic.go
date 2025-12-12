// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package delegate

import (
	"context"

	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelegateWithdrawGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelegateWithdrawGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelegateWithdrawGetLogic {
	return &DelegateWithdrawGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelegateWithdrawGetLogic) DelegateWithdrawGet(req *types.DelegateWithdrawGetReq) (resp *types.DelegateWithdrawItem, err error) {
	// todo: add your logic here and delete this line

	return
}
