package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"jxChainWebSvc/internal/logic"
	"jxChainWebSvc/internal/svc"
	"jxChainWebSvc/internal/types"
)

func JxChainDataDeleteHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteData
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewJxChainDataDeleteLogic(r.Context(), ctx)
		resp, err := l.JxChainDataDelete(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
