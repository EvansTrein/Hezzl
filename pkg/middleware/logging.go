package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func HandlerLog(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			op := "http request"

			var builder strings.Builder
			builder.WriteString("\n")
			for key, values := range r.Header {
				builder.WriteString("   ")
				builder.WriteString(fmt.Sprintf("%s: %s\n", key, values))
			}

			log := log.With(
				slog.String("operation", op),
				slog.String("path", r.URL.Path),
				slog.String("HTTP Method", r.Method),
				slog.String("ip-address", r.RemoteAddr),
				slog.String("host", r.Host),
			)
			log.Debug("request received", "headers", builder.String())

			next.ServeHTTP(w, r)

			log.Info("request completed", "duration", time.Since(start))
		})
	}
}
