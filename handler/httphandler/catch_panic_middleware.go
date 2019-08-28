package httphandler

import (
	"net/http"

	"go.uber.org/zap"
)

// CatchPanicMiddleware describe catch panic middleware
type CatchPanicMiddleware struct {
	logger *zap.Logger
}

// NewCatchPanicMiddleware creates new middleware
func NewCatchPanicMiddleware(logger *zap.Logger) *CatchPanicMiddleware {
	return &CatchPanicMiddleware{
		logger: logger,
	}
}

// Middleware implements panic catcher middleware
func (m *CatchPanicMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if recErr := recover(); recErr != nil {
				m.logger.With(
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
					zap.String("mode", "panic_log"),
					zap.Any("panic", recErr),
					zap.Stack("stack"),
				).Error("Panic happened in HTTP handler")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
