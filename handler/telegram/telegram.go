package telegram

import (
	"context"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"gitlab.thedatron.ru/anton/message_gateway/command"
	"gitlab.thedatron.ru/anton/message_gateway/handler"
	"gitlab.thedatron.ru/anton/message_gateway/service"
	"go.uber.org/zap"
)

const (
	updatesTimeout = 60
	errorOccurred  = "Error occurred during processing"
)

type telegram struct {
	bot             *tgbotapi.BotAPI
	logger          *zap.Logger
	commands        map[string]command.Command
	shutdownChannel chan interface{}
}

// NewClient creates new instance of telegram
func NewClient(vars service.Vars,
	client *http.Client, logger *zap.Logger, commands []command.Command) (handler.Sender, handler.Receiver) {

	_ = tgbotapi.SetLogger(&botLogger{})
	bot, err := tgbotapi.NewBotAPIWithClient(vars.TelegramToken, client)
	if err != nil {
		logger.Fatal("Failed to initialize Telegram", zap.Error(err))
	}
	logger.Info("Telegram bot initialized")

	telegram := &telegram{
		bot:             bot,
		logger:          logger,
		commands:        createMapping(commands),
		shutdownChannel: make(chan interface{}),
	}
	return telegram, telegram
}

func createMapping(commands []command.Command) map[string]command.Command {
	mapping := make(map[string]command.Command)
	for _, cmd := range commands {
		mapping[cmd.Command()] = cmd
	}
	return mapping
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

func (t *telegram) Receive(ctx context.Context) {
	go func() {
		config := tgbotapi.NewUpdate(0)
		config.Timeout = updatesTimeout
		updates, _ := t.bot.GetUpdatesChan(config)

		for {
			select {
			case <-t.shutdownChannel:
				return
			case update := <-updates:
				message := update.Message
				if message == nil {
					message = update.ChannelPost
					if message == nil {
						continue
					}
				}

				if message.IsCommand() {
					cmd := t.commands[message.Command()]
					if cmd != nil {
						go func(ctx context.Context, chatID int64) { // Run command processing async
							text, err := cmd.Reply(ctx)
							if err != nil {
								t.logger.Error("Failed to process reply", zap.Error(err))
								text = errorOccurred
							}

							reply := tgbotapi.NewMessage(chatID, text)
							if _, err := t.bot.Send(&reply); err != nil {
								t.logger.Error("Failed to send reply", zap.Error(err))
							}
						}(ctx, message.Chat.ID)
					}
				}
			}
		}
	}()
}

func (t *telegram) Stop() {
	t.bot.StopReceivingUpdates()
	close(t.shutdownChannel)

	t.logger.Info("Telegram bot stopped")
}
