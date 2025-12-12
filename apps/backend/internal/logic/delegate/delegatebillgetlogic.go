// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package delegate

import (
	"context"

	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelegateBillGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelegateBillGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelegateBillGetLogic {
	return &DelegateBillGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelegateBillGetLogic) DelegateBillGet(req *types.DelegateBillGetReq) (resp *types.DelegateBillItem, err error) {
	// todo: add your logic here and delete this line

	return
}
