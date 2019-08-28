package handler

// MessageTransport type of transport to send message
type MessageTransport int64

const (
	SMSMessageTransport MessageTransport = iota + 1
	TelegramMessageTransport
	EmailMessageTransport
)

// SendMessage declare send message model
type SendMessage struct {
	Transport MessageTransport `json:"transport"`
	Target    string           `json:"target"`
	Text      string           `json:"text"`
}
