// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package delegate

import (
	"context"

	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type DelegateOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelegateOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelegateOrderListLogic {
	return &DelegateOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelegateOrderListLogic) DelegateOrderList(req *types.DelegateOrderListReq) (resp *types.DelegateOrderListResp, err error) {
	// todo: add your logic here and delete this line
	orders, total, err := l.svcCtx.OrderService.FindDelegateOrderList(req.TransactionId, req.FromBase58, req.ToBase58, req.Typo, req.Status, req.Start, req.End, req.Page, req.Limit)
	if err != nil {
		logx.Errorf("delegate order list find failed, err:%v", err)
		return
	}

	items := make([]*types.DelegateOrderItem, len(orders))
	copier.Copy(&items, &orders)
	resp = &types.DelegateOrderListResp{Items: items, Total: total}
	return
}
