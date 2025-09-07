package sloglog

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RequestLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ctx := r.Context()
			log := FromContext(ctx).With(slog.String("uri", r.RequestURI), slog.String("method", r.Method))

			r = r.WithContext(ToContext(ctx, log))

			scheme := "http"

			if r.TLS != nil {
				scheme = "https"
			}

			t1 := time.Now()
			defer func() {
				log.InfoContext(ctx, "request",
					slog.Any("status", ww.Status()),
					slog.Any("bytes_written", ww.BytesWritten()),
					slog.String("time_since", strconv.Itoa(int(time.Since(t1).Milliseconds()))+"ms"),
					slog.String("host", r.Host),
					slog.String("proto", r.Proto),
					slog.String("url", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)),
					slog.Any("params", chi.RouteContext(ctx).URLParams),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
