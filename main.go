package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Toshik1978/message_gateway/handler/telegram"
	"github.com/Toshik1978/message_gateway/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	GitVersion = "undefined"

	interruptCh chan os.Signal
)

func main() {
	logger := initializeLogger()
	defer func() {
		if recErr := recover(); recErr != nil {
			// Log error
			logger.Error("Panic in main", zap.Any("panic", recErr))
		}
	}()

	interruptCh = make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Start service", zap.String("git_version", GitVersion))

	vars := service.LoadConfig(logger)
	httpClient := service.NewHTTPClient(vars, logger)
	_ = telegram.NewTelegram(vars, httpClient, logger)
}

// initializeLogger initializes logger
func initializeLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		log.Fatal("Initialize logger failed", err)
	}
	return logger
}
