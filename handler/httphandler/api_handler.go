package httphandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Toshik1978/message_gateway/handler"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// APIHandler declare handler API requests
type APIHandler struct {
	logger  *zap.Logger
	version string

	senders map[handler.MessageTransport]handler.Sender
}

// NewAPIHandler creates new API handler
func NewAPIHandler(senders map[handler.MessageTransport]handler.Sender, logger *zap.Logger, version string) *APIHandler {
	return &APIHandler{
		logger:  logger,
		version: version,
		senders: senders,
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

		errs := make([]error, 0)
		sent := false
		for _, transport := range message.Transports {
			if sender, ok := h.senders[transport]; ok {
				err := sender.Send(ctx, message.Target, message.Text)
				if err == nil {
					sent = true
					break
				} else {
					errs = append(errs, err)
				}
			}
		}

		if !sent {
			errText := ""
			for _, err := range errs {
				errText += err.Error() + "\n"
			}
			_ = h.fail(w, errors.Errorf("failed to send message: %s", errText), "SendHandler")
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
