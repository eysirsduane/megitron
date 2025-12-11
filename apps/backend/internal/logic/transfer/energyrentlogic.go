// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package transfer

import (
	"context"

	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnergyRentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnergyRentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnergyRentLogic {
	return &EnergyRentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnergyRentLogic) EnergyRent(req *types.ResourceDelegateReq) (resp *types.ResourceDelegateResp, err error) {
	// todo: add your logic here and delete this line

	return
}
