package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Toshik1978/message_gateway/handler"
	"go.uber.org/zap"
)

// APIHandler declare handler API requests
type APIHandler struct {
	logger  *zap.Logger
	version string

	telegramClient handler.Sender
}

// NewAPIHandler creates new API handler
func NewAPIHandler(telegramClient handler.Sender, logger *zap.Logger, version string) *APIHandler {
	return &APIHandler{
		logger:         logger,
		version:        version,
		telegramClient: telegramClient,
	}
}

// serviceStatus declare status structure
type serviceStatus struct {
	IsAlive bool      `json:"is_alive"`
	Version string    `json:"version"`
	Date    time.Time `json:"date"`
}

// ServiceStatusHandler sends status of the service and it's version
func (h *APIHandler) ServiceStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.writeResponse(w, serviceStatus{
			IsAlive: true,
			Version: h.version,
			Date:    time.Now(),
		})
	})
}

// ServiceStatusHandler sends status of the service and it's version
func (h *APIHandler) SendHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		var message handler.SendMessage
		decoder := json.NewDecoder(r.Body)
		if h.fail(w, decoder.Decode(&message), "SendHandler") {
			return
		}

		var err error
		switch message.Transport {
		case handler.SMSMessageTransport:
			err = errors.New("unsupported sms transport")
		case handler.TelegramMessageTransport:
			err = h.telegramClient.Send(ctx, message.Target, message.Text)
		case handler.EmailMessageTransport:
			err = errors.New("unsupported email transport")
		default:
			err = errors.New("failed to detect valid message transport")
		}

		if h.fail(w, err, "SendHandler") {
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

// fail fails request
func (h *APIHandler) fail(w http.ResponseWriter, err error, method string) bool {
	if err != nil {
		h.logger.Error("Failed to handle "+method, zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

// writeResponse write response
func (h *APIHandler) writeResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	payload, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal HTTP response", zap.Error(err), zap.Any("response", response))
		return
	}
	_, err = w.Write(payload)
	if err != nil {
		h.logger.Error("Failed to write HTTP response", zap.Error(err), zap.Any("response", response))
		return
	}
}
