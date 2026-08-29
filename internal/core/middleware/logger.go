package core_middleware

import (
	"net/http"
	"time"
	core_logger "workout_app/internal/core/logger"

	"go.uber.org/zap"
)

func Logger(log *core_logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := NewResponseWriter(w)

			start := time.Now()

			next.ServeHTTP(rw, r)

			log.Debug(
				"HTTP request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status_code", rw.statusCode),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}
