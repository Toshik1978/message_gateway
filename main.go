package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Toshik1978/message_gateway/handler"
	"github.com/Toshik1978/message_gateway/handler/httphandler"
	"github.com/Toshik1978/message_gateway/handler/telegram"
	"github.com/Toshik1978/message_gateway/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	GitVersion = "undefined"

	interruptCh         chan os.Signal
	httpShutdownTimeout = 5 * time.Second
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
	telegramClient := telegram.NewTelegram(vars, httpClient, logger)
	server := initializeHTTP(vars, telegramClient, logger)

	waitShutdown(logger, server)
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

// initializeHTTP initializes HTTP server
func initializeHTTP(vars service.Vars, telegramClient handler.Sender, logger *zap.Logger) *http.Server {
	apiHandler := httphandler.NewAPIHandler(telegramClient, logger, GitVersion)
	r := mux.NewRouter()
	r.Use(httphandler.NewCatchPanicMiddleware(logger).Middleware)
	r.Use(httphandler.NewAccessLogMiddleware(logger).Middleware)

	r.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)
	r.PathPrefix("/status").Handler(apiHandler.ServiceStatusHandler()).Methods("GET")

	route := r.PathPrefix("/api/v1").Subrouter()
	route.Handle("/send", apiHandler.SendHandler()).Methods("POST")

	server := &http.Server{
		Addr:    vars.HTTPAddress,
		Handler: r,
	}

	go func() {
		logger.Info("HTTP server initializing",
			zap.String("http_addr", vars.HTTPAddress))
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Info("HTTP server shutdown failed", zap.Error(err))
		} else {
			logger.Info("HTTP server shutdown")
		}
	}()

	return server
}

func waitShutdown(logger *zap.Logger, server *http.Server) {
	// Wait for interrupt
	<-interruptCh

	ctx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Info("Failed to graceful shutdown server", zap.Error(err))
	}
	server = nil

	logger.Info("Stop service", zap.String("git_version", GitVersion))
}
