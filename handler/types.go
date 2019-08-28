package handler

// MessageTransport type of transport to send message
type MessageTransport int64

const (
	SMSMessageTransport MessageTransport = iota + 1
	TelegramMessageTransport
	EmailMessageTransport
)
