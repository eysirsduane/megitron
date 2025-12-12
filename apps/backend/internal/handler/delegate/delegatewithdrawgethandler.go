// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package delegate

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"megitron/apps/backend/internal/logic/delegate"
	"megitron/apps/backend/internal/svc"
	"megitron/apps/backend/internal/types"
)

func DelegateWithdrawGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DelegateWithdrawGetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := delegate.NewDelegateWithdrawGetLogic(r.Context(), svcCtx)
		resp, err := l.DelegateWithdrawGet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
