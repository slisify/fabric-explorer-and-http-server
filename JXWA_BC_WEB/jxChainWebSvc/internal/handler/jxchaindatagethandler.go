package handler

import (
	"net/http"

	"jxChainWebSvc/internal/logic"
	"jxChainWebSvc/internal/svc"
	"jxChainWebSvc/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func JxChainDataGetHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryData
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewJxChainDataGetLogic(r.Context(), ctx)
		resp, err := l.JxChainDataGet(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
