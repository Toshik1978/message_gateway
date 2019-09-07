package speedtest

import (
	"context"
	"fmt"

	"github.com/kylegrantlucas/speedtest"
	"github.com/pkg/errors"
	"gitlab.thedatron.ru/anton/message_gateway/command"
	"gitlab.thedatron.ru/anton/message_gateway/service"
	"go.uber.org/zap"
)

type speedCommand struct {
	logger *zap.Logger
	client *speedtest.Client
}

// NewCommand creates command to process speed
func NewCommand(vars service.Vars, logger *zap.Logger) command.Command {
	client, err := speedtest.NewDefaultClient()
	if err != nil {
		logger.Fatal("Failed to initialize speed")
	}

	logger.Info("Speed command started")
	return &speedCommand{
		logger: logger,
		client: client,
	}
}

func (c *speedCommand) Command() string {
	return "speed"
}

func (c *speedCommand) Reply(ctx context.Context) (string, error) {
	server, err := c.client.GetServer("")
	if err != nil {
		return "", errors.Wrap(err, "failed to get speedtest client")
	}
	download, err := c.client.Download(server)
	if err != nil {
		return "", errors.Wrap(err, "failed to test download")
	}
	upload, err := c.client.Upload(server)
	if err != nil {
		return "", errors.Wrap(err, "failed to test upload")
	}
	return fmt.Sprintf("Download: %.2f\nUpload: %.2f", download, upload), nil
}
