package telegram

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Toshik1978/message_gateway/handler"
	"github.com/Toshik1978/message_gateway/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type telegram struct {
	bot    *tgbotapi.BotAPI
	logger *zap.Logger
}

// NewTelegram creates new instance of telegram
func NewTelegram(vars service.Vars, client *http.Client, logger *zap.Logger) handler.Sender {
	bot, err := tgbotapi.NewBotAPIWithClient(vars.TelegramToken, client)
	if err != nil {
		logger.Fatal("Failed to initialize Telegram", zap.Error(err))
	}
	logger.Info("Telegram bot initialized")
	return &telegram{
		bot:    bot,
		logger: logger,
	}
}

func (t *telegram) Send(ctx context.Context, target string, text string) error {
	t.logger.Info("Send telegram message", zap.String("to", target))

	// target is channel id
	channelID, err := strconv.ParseInt(target, 10, 64)
	if err != nil {
		return errors.Wrap(err, "failed to get channel id")
	}

	message := tgbotapi.NewMessage(channelID, text)
	_, err = t.bot.Send(&message)
	return errors.Wrap(err, "failed to send message via telegram")
}
