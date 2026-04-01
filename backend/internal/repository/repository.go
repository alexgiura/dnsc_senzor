package repository

type Repository struct {
	NetworkAlerts NetworkAlertFileRepository
}

func NewRepository(networkAlertsStoragePath string) *Repository {
	return &Repository{
		NetworkAlerts: NewNetworkAlertFileRepository(networkAlertsStoragePath),
	}
}
