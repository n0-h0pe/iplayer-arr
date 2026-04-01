package download

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"time"

	"github.com/GiteaLN/iplayer-arr/internal/store"
)

type Manager struct {
	store       *store.Store
	downloadDir string
	maxWorkers  int
}

func NewManager(st *store.Store, downloadDir string, maxWorkers int) *Manager {
	return &Manager{
		store:       st,
		downloadDir: downloadDir,
		maxWorkers:  maxWorkers,
	}
}

func (m *Manager) Enqueue(pid, quality, title, category string) (string, error) {
	existing, _ := m.store.FindDownloadByPIDQuality(pid, quality)
	if existing != nil {
		return existing.ID, nil
	}

	hist, _ := m.store.FindHistoryByPIDQuality(pid, quality)
	if hist != nil {
		return hist.ID, nil
	}

	id := generateNzoID()
	outputDir := filepath.Join(m.downloadDir, title)

	dl := &store.Download{
		ID:        id,
		PID:       pid,
		Quality:   quality,
		Title:     title,
		Category:  category,
		Status:    store.StatusPending,
		OutputDir: outputDir,
		CreatedAt: time.Now(),
	}

	if err := m.store.PutDownload(dl); err != nil {
		return "", fmt.Errorf("store download: %w", err)
	}

	return id, nil
}

func (m *Manager) CancelDownload(nzoID string) error {
	m.store.DeleteDownload(nzoID)
	return nil
}

func (m *Manager) StartDownload(pid, quality, title, category string) (string, error) {
	return m.Enqueue(pid, quality, title, category)
}

func generateNzoID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return "iparr_" + hex.EncodeToString(b)
}
