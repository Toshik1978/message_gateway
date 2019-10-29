package ups

import (
	"context"
	"fmt"

	"github.com/Toshik1978/message_gateway/command"
	"github.com/Toshik1978/message_gateway/service"
	"github.com/pkg/errors"
	nut "github.com/robbiet480/go.nut"
	"go.uber.org/zap"
)

const (
	voltage = "input.voltage"
)

type upsCommand struct {
	logger  *zap.Logger
	address string
	name    string
	login   string
	pass    string
}

// NewCommand creates command to process UPS
func NewCommand(vars service.Vars, logger *zap.Logger) command.Command {
	logger.Info("UPS command started")
	return &upsCommand{
		logger:  logger,
		address: vars.UpsAddress,
		name:    vars.UpsName,
		login:   vars.UpsLogin,
		pass:    vars.UpsPass,
	}
}

func (c *upsCommand) Command() string {
	return "power"
}

func (c *upsCommand) Reply(ctx context.Context) (string, error) {
	ups, client, err := c.connect()
	if err != nil {
		return "", errors.Wrap(err, "failed to connect to ups")
	}
	defer c.disconnect(client)

	vars, err := ups.GetVariables()
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

func (c *upsCommand) connect() (*nut.UPS, *nut.Client, error) {
	client, err := nut.Connect(c.address)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to connect to nut")
	}
	if ok, err := client.Authenticate(c.login, c.pass); !ok || err != nil {
		return nil, nil, errors.Wrap(err, "failed to auth to nut")
	}
	ups, err := nut.NewUPS(c.name, &client)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to initialize ups")
	}
	return &ups, &client, nil
}

func (c *upsCommand) disconnect(client *nut.Client) {
	_, _ = client.Disconnect()
}
