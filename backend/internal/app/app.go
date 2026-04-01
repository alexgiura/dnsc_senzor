package app

import (
	"context"
	"fmt"

	"senzor/internal/config"
	"senzor/internal/middleware"
	"senzor/internal/repository"
	"senzor/internal/routes"
	"senzor/internal/server"
	"senzor/internal/services"
	"senzor/internal/utils"
)

type App struct {
	server *server.Server
	logger *utils.Logger
}

func NewApp(cfg *config.Config, logger *utils.Logger) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	repo := repository.NewRepository(cfg.NetworkAlerts.StoragePath)
	appServices := services.NewAppServices(repo)

	router := routes.RegisterRoutes(appServices)

	handlerWithMiddleware := middleware.CorsMiddleware(router)

	srv, err := server.NewServer(cfg.AppSettings.ServerPort, handlerWithMiddleware)
	if err != nil {
		return nil, fmt.Errorf("create server: %w", err)
	}

	return &App{
		server: srv,
		logger: logger,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	if app.server == nil {
		return fmt.Errorf("server is nil")
	}

	if app.logger != nil {
		app.logger.Info("APP", "startup", 0, 0, "starting HTTP server")
	}

	if err := app.server.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}

	return nil
}

func (app *App) Shutdown(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("shutdown context is nil")
	}

	if app.logger != nil {
		app.logger.Info("APP", "shutdown", 0, 0, "shutting down application")
	}

	var shutdownErr error

	if app.server != nil {
		if err := app.server.Shutdown(ctx); err != nil {
			if app.logger != nil {
				app.logger.Error("APP", "shutdown", 0, 0, fmt.Sprintf("error shutting down server: %v", err))
			}
			shutdownErr = fmt.Errorf("shutdown server: %w", err)
		} else if app.logger != nil {
			app.logger.Info("APP", "shutdown", 0, 0, "server stopped successfully")
		}
	}

	return shutdownErr
}
