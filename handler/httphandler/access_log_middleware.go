package httphandler

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// AccessLogMiddleware middleware for access log
type AccessLogMiddleware struct {
	logger *zap.Logger
}

// NewAccessLogMiddleware creates new middleware
func NewAccessLogMiddleware(logger *zap.Logger) *AccessLogMiddleware {
	return &AccessLogMiddleware{
		logger: logger,
	}
}

// Middleware implements access log middleware
func (m *AccessLogMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		m.logger.With(
			zap.Duration("duration", time.Since(start)),
			zap.String("url", r.URL.String()),
			zap.String("remote addr", r.RemoteAddr),
			zap.String("method", r.Method),
			zap.String("real ip", r.Header.Get("X-Real-Ip")),
			zap.String("user agent", r.UserAgent()),
			zap.String("cm_user", r.Header.Get("CM-User")),
			zap.String("cm_session", r.Header.Get("CM-Session")),
			zap.String("cm_device", r.Header.Get("CM-Device")),
			zap.String("cm_auth_token", r.Header.Get("CM-Auth-Token")),
			zap.String("app_version", r.Header.Get("CM-App-Version")),
			zap.String("mode", "access_log"),
		).Info("Handled request")
	})
}
