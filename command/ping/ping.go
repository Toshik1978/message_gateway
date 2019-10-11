package ping

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sparrc/go-ping"
	"gitlab.thedatron.ru/anton/message_gateway/command"
	"gitlab.thedatron.ru/anton/message_gateway/service"
	"go.uber.org/zap"
)

type pingCommand struct {
	logger *zap.Logger
	hosts  []string
}

// NewCommand creates command to process ping
func NewCommand(vars service.Vars, logger *zap.Logger) command.Command {
	logger.Info("Ping command started")

	return &pingCommand{
		logger: logger,
		hosts:  strings.Split(vars.PingHosts, ","),
	}
}

func (c *pingCommand) Command() string {
	return "ping"
}

func (c *pingCommand) Reply(ctx context.Context) (string, error) {
	text := "PONG\n"
	for _, host := range c.hosts {
		pinger, err := ping.NewPinger(host)
		if err != nil {
			return "", errors.Wrap(err, "failed to ping")
		}
		pinger.Count = 3
		pinger.Timeout = 5 * time.Second
		pinger.Run()
		statistics := pinger.Statistics()
		if statistics.PacketsSent == statistics.PacketsRecv {
			text += host + ": OK\n"
		} else {
			text += host + ": FAIL\n"
		}
	}
	return text, nil
}
