package handler

import (
	"net/http"

	"jxChainWebSvc/internal/logic"
	"jxChainWebSvc/internal/svc"
	"jxChainWebSvc/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func JxChainDataPutHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OnChainDataPut
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewJxChainDataPutLogic(r.Context(), ctx)
		resp, err := l.JxChainDataPut(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
