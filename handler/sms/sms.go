package sms

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Toshik1978/message_gateway/handler"
	"github.com/Toshik1978/message_gateway/service"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	tmpPrefix = "sms_"
)

type sms struct {
	incomingPath string
	outgoingPath string
	tempPath     string

	logger *zap.Logger
}

// NewSMS creates new instance of sms
func NewSMS(vars service.Vars, logger *zap.Logger) handler.Sender {
	return &sms{
		incomingPath: vars.SMSIncomingPath,
		outgoingPath: vars.SMSOutgoingPath,
		tempPath:     vars.SMSTempPath,
		logger:       logger,
	}
}

func (e *sms) Send(ctx context.Context, target string, text string) error {
	e.logger.Info("Send SMS", zap.String("to", target))

	message :=
		"To: " + target + "\n" +
			"Alphabet: ISO\n" +
			"UDH: false\n\n" +
			text + "\n"

	fileName, err := e.createTempFile(message)
	if err != nil {
		return errors.Wrap(err, "failed to create tmp file")
	}
	return errors.Wrap(e.moveFile(fileName), "failed to send message via email")
}

// createTempFile creates temp file with message
func (e *sms) createTempFile(message string) (string, error) {
	file, err := ioutil.TempFile(e.tempPath, tmpPrefix)
	if err != nil {
		return "", errors.Wrap(err, "failed to create tmp file")
	}

	if _, err := file.Write([]byte(message)); err != nil {
		_ = file.Close()
		return "", errors.Wrap(err, "failed to store message in file")
	}

	if err := file.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close file")
	}

	return file.Name(), nil
}

// moveFile moves temp file to outgoing path
func (e *sms) moveFile(fileName string) error {
	_, file := filepath.Split(fileName)
	if err := os.Rename(fileName, filepath.Join(e.outgoingPath, file)); err != nil {
		return errors.Wrap(err, "failed to move file")
	}
	return nil
}
