package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"senzor/internal/models"
)

type NetworkAlertFileRepository interface {
	Append(ctx context.Context, alert models.NetworkAlert) error
}

type networkAlertFileRepository struct {
	path string
	mu   sync.Mutex
}

func NewNetworkAlertFileRepository(storagePath string) NetworkAlertFileRepository {
	return &networkAlertFileRepository{path: storagePath}
}

func (r *networkAlertFileRepository) Append(_ context.Context, alert models.NetworkAlert) error {
	if r.path == "" {
		return fmt.Errorf("network alerts storage path is empty")
	}

	line, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("marshal alert: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return fmt.Errorf("create storage directory: %w", err)
	}

	f, err := os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open storage file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(line); err != nil {
		return fmt.Errorf("write alert: %w", err)
	}
	if _, err := f.Write([]byte("\n")); err != nil {
		return fmt.Errorf("write newline: %w", err)
	}

	return nil
}
