// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package exchange

import (
	"context"

	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExchangeOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExchangeOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeOrderListLogic {
	return &ExchangeOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExchangeOrderListLogic) ExchangeOrderList(req *types.ExchangeOrderListReq) (resp *types.ExchangeOrderListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
