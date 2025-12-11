// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package account

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megitron/apps/backend/internal/logic/account"
	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"
)

func TronAccountCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TronAccountCreateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := account.NewTronAccountCreateLogic(r.Context(), svcCtx)
		resp, err := l.TronAccountCreate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
