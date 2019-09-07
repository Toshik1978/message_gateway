package ups

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	nut "github.com/robbiet480/go.nut"
	"gitlab.thedatron.ru/anton/message_gateway/command"
	"gitlab.thedatron.ru/anton/message_gateway/service"
	"go.uber.org/zap"
)

const (
	voltage = "input.voltage"
)

type upsCommand struct {
	client nut.Client
	ups    nut.UPS
}

// NewCommand creates command to process UPS
func NewCommand(vars service.Vars, logger *zap.Logger) command.Command {
	client, err := nut.Connect(vars.UpsAddress)
	if err != nil {
		logger.Fatal("Failed to initialize NUT", zap.Error(err))
	}
	if ok, err := client.Authenticate(vars.UpsLogin, vars.UpsPass); !ok || err != nil {
		logger.Fatal("Failed to authenticate in NUT", zap.Error(err))
	}
	ups, err := nut.NewUPS(vars.UpsName, &client)
	if err != nil {
		logger.Fatal("Failed to initialize UPS", zap.Error(err))
	}

	logger.Info("UPS command started")
	return &upsCommand{
		client: client,
		ups:    ups,
	}
}

func (c *upsCommand) Command() string {
	return "power"
}

func (c *upsCommand) Reply(ctx context.Context) (string, error) {
	vars, err := c.ups.GetVariables()
	if err != nil {
		return "", errors.Wrap(err, "failed to get nut variables")
	}

	for _, variable := range vars {
		if variable.Name == voltage {
			value := variable.Value.(float64)
			return fmt.Sprintf("Voltage: %.1f", value), nil
		}
	}
	return "Voltage: unknown", nil
}
