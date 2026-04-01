package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"senzor/internal/app"
	"senzor/internal/config"
	"senzor/internal/utils"
)

const shutdownTimeout = 10 * time.Second

func main() {
	os.Exit(run())
}

func run() int {
	logger := utils.GetLogger("senzor")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("APP", "startup", 0, 0, "failed to load config: "+err.Error())
		return 1
	}

	application, err := app.NewApp(cfg, logger)
	if err != nil {
		logger.Error("APP", "startup", 0, 0, "failed to initialize app: "+err.Error())
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		errCh <- application.Run(ctx)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			logger.Error("APP", "runtime", 0, 0, "application stopped with error: "+err.Error())
			return 1
		}
		logger.Info("APP", "runtime", 0, 0, "application stopped gracefully")
		return 0

	case <-ctx.Done():
		logger.Info("APP", "shutdown", 0, 0, "shutdown initiated")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := application.Shutdown(shutdownCtx); err != nil {
			logger.Error("APP", "shutdown", 0, 0, "shutdown error: "+err.Error())
			return 1
		}

		select {
		case err := <-errCh:
			if err != nil {
				logger.Error("APP", "shutdown", 0, 0, "application stopped with error: "+err.Error())
				return 1
			}
			logger.Info("APP", "shutdown", 0, 0, "application stopped gracefully")
			return 0
		case <-shutdownCtx.Done():
			logger.Error("APP", "shutdown", 0, 0, "timed out waiting for application to stop")
			return 1
		}
	}
}
