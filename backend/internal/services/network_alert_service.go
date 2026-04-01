package services

import (
	"context"
	"fmt"
	"strings"

	apperrors "senzor/internal/errors"
	"senzor/internal/models"
	"senzor/internal/repository"
)

type NetworkAlertService interface {
	Ingest(ctx context.Context, alert *models.NetworkAlert) error
}

type networkAlertService struct {
	repo repository.NetworkAlertFileRepository
}

func NewNetworkAlertService(repo repository.NetworkAlertFileRepository) NetworkAlertService {
	return &networkAlertService{repo: repo}
}

func (s *networkAlertService) Ingest(ctx context.Context, alert *models.NetworkAlert) error {
	if alert == nil {
		return fmt.Errorf("nil alert: %w", apperrors.ErrValidation)
	}

	if err := validateNetworkAlert(*alert); err != nil {
		return fmt.Errorf("%w: %s", apperrors.ErrValidation, err)
	}

	return s.repo.Append(ctx, *alert)
}

func validateNetworkAlert(a models.NetworkAlert) error {
	if strings.TrimSpace(a.AgentID) == "" {
		return fmt.Errorf("agent_id is required")
	}
	if a.Event.Protocol == "" {
		return fmt.Errorf("event.protocol is required")
	}
	if strings.TrimSpace(a.Event.SrcIP) == "" || strings.TrimSpace(a.Event.DstIP) == "" {
		return fmt.Errorf("event.src_ip and event.dst_ip are required")
	}
	if strings.TrimSpace(a.Event.WatchlistMatch) == "" {
		return fmt.Errorf("event.watchlist_match is required")
	}
	if strings.TrimSpace(a.Event.Direction) == "" {
		return fmt.Errorf("event.direction is required")
	}
	return nil
}
