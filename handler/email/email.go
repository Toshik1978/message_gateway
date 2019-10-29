package email

import (
	"context"
	"net/smtp"

	"github.com/Toshik1978/message_gateway/handler"
	"github.com/Toshik1978/message_gateway/service"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type email struct {
	name     string
	subject  string
	smtp     string
	smtpPort string
	login    string
	pass     string

	logger *zap.Logger
}

// NewClient creates new instance of email
func NewClient(vars service.Vars, logger *zap.Logger) handler.Sender {
	logger.Info("Email initialized")

	return &email{
		name:     vars.EmailName,
		subject:  vars.EmailSubject,
		smtp:     vars.EmailSMTP,
		smtpPort: vars.EmailPort,
		login:    vars.EmailLogin,
		pass:     vars.EmailPass,
		logger:   logger,
	}
}

func (e *email) Send(ctx context.Context, target string, text string) error {
	e.logger.Info("Send email", zap.String("to", target))

	message :=
		"From: " + e.name + "\n" +
			"To: " + target + "\n" +
			"Subject: " + e.subject + "\n" +
			text

	err := smtp.SendMail(e.smtp+":"+e.smtpPort,
		smtp.PlainAuth("", e.login, e.pass, e.smtp),
		e.name, []string{target}, []byte(message))
	return errors.Wrap(err, "failed to send message via email")
}
