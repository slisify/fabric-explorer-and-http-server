package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/sysx"
	"github.com/tal-tech/go-zero/core/trace"
)

// TracingHandler returns a middleware that traces the request.
func TracingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		carrier, err := trace.Extract(trace.HttpFormat, r.Header)
		// ErrInvalidCarrier means no trace id was set in http header
		if err != nil && err != trace.ErrInvalidCarrier {
			logx.Error(err)
		}

		ctx, span := trace.StartServerSpan(r.Context(), carrier, sysx.Hostname(), r.RequestURI)
		defer span.Finish()
		r = r.WithContext(ctx)

		// conveniently tracking error messages
		w.Header().Set(trace.TraceIdKey, span.TraceId())
		next.ServeHTTP(w, r)
	})
}
