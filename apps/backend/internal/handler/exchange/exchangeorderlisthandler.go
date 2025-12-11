// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package exchange

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megitron/apps/backend/internal/logic/exchange"
	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"
)

func ExchangeOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ExchangeOrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := exchange.NewExchangeOrderListLogic(r.Context(), svcCtx)
		resp, err := l.ExchangeOrderList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
