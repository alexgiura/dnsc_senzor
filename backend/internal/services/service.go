package services

import (
	"senzor/internal/repository"
)

type AppServices struct {
	NetworkAlert NetworkAlertService
}

func NewAppServices(repos *repository.Repository) *AppServices {
	return &AppServices{
		NetworkAlert: NewNetworkAlertService(repos.NetworkAlerts),
	}
}
