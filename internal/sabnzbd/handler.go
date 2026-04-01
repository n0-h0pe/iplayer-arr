package sabnzbd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GiteaLN/iplayer-arr/internal/store"
)

type DownloadStarter interface {
	StartDownload(pid, quality, title, category string) (string, error)
	CancelDownload(nzoID string) error
}

type Handler struct {
	store   *store.Store
	starter DownloadStarter
}

func NewHandler(st *store.Store, starter DownloadStarter) *Handler {
	return &Handler{store: st, starter: starter}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")

	switch mode {
	case "version":
		w.Write([]byte("4.0.0"))
		return
	case "get_cats":
		writeJSON(w, map[string]interface{}{"categories": []string{"sonarr", "tv", "manual"}})
		return
	}

	// all other modes require auth
	apiKey := r.URL.Query().Get("apikey")
	storedKey, _ := h.store.GetConfig("api_key")
	if apiKey != storedKey {
		writeJSON(w, map[string]interface{}{
			"status": false,
			"error":  "API Key Incorrect",
		})
		return
	}

	switch mode {
	case "queue":
		h.handleQueue(w, r)
	case "history":
		h.handleHistory(w, r)
	case "addurl", "addfile":
		h.handleAdd(w, r)
	default:
		writeJSON(w, map[string]interface{}{"status": false, "error": "unknown mode"})
	}
}

func (h *Handler) handleQueue(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "delete" {
		value := r.URL.Query().Get("value")
		if h.starter != nil {
			h.starter.CancelDownload(value)
		}
		h.store.DeleteDownload(value)
		writeJSON(w, map[string]interface{}{"status": true})
		return
	}

	downloads, _ := h.store.ListDownloads()
	var slots []map[string]interface{}
	for _, dl := range downloads {
		if dl.Status == store.StatusCompleted || dl.Status == store.StatusFailed {
			continue
		}

		status := "Queued"
		switch dl.Status {
		case store.StatusDownloading, store.StatusConverting:
			status = "Downloading"
		case store.StatusResolving:
			status = "Queued"
		}

		mbTotal := float64(dl.Size) / 1024 / 1024
		mbLeft := mbTotal * (1 - dl.Progress/100)

		slots = append(slots, map[string]interface{}{
			"nzo_id":     dl.ID,
			"filename":   dl.Title,
			"status":     status,
			"percentage": fmt.Sprintf("%.0f", dl.Progress),
			"mb":         fmt.Sprintf("%.2f", mbTotal),
			"mbleft":     fmt.Sprintf("%.2f", mbLeft),
			"timeleft":   "unknown",
			"cat":        dl.Category,
		})
	}

	if slots == nil {
		slots = []map[string]interface{}{}
	}

	writeJSON(w, map[string]interface{}{
		"queue": map[string]interface{}{
			"slots": slots,
		},
	})
}

func (h *Handler) handleHistory(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "delete" {
		value := r.URL.Query().Get("value")
		h.store.DeleteHistory(value)
		writeJSON(w, map[string]interface{}{"status": true})
		return
	}

	history, _ := h.store.ListHistory()
	var slots []map[string]interface{}
	for _, dl := range history {
		status := "Completed"
		if dl.Status == store.StatusFailed {
			status = "Failed"
		}
		slots = append(slots, map[string]interface{}{
			"nzo_id":  dl.ID,
			"name":    dl.Title,
			"status":  status,
			"storage": dl.OutputDir,
			"bytes":   dl.Size,
			"cat":     dl.Category,
		})
	}

	if slots == nil {
		slots = []map[string]interface{}{}
	}

	writeJSON(w, map[string]interface{}{
		"history": map[string]interface{}{
			"slots": slots,
		},
	})
}

func (h *Handler) handleAdd(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]interface{}{
		"status":  true,
		"nzo_ids": []string{"iparr_placeholder"},
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
