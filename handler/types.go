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
	Messages []struct {
		Transport MessageTransport `json:"transport"`
		Target    string           `json:"target"`
	} `json:"messages"`
	Text string `json:"text"`
}
