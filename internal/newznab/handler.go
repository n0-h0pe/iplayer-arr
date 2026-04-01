package newznab

import (
	"net/http"

	"github.com/GiteaLN/iplayer-arr/internal/bbc"
	"github.com/GiteaLN/iplayer-arr/internal/store"
)

type Handler struct {
	ibl   *bbc.IBL
	store *store.Store
	ms    *bbc.MediaSelector
}

func NewHandler(ibl *bbc.IBL, st *store.Store, ms *bbc.MediaSelector) *Handler {
	return &Handler{ibl: ibl, store: st, ms: ms}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")

	switch t {
	case "caps":
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(capsXML()))
	case "search":
		h.handleSearch(w, r)
	case "tvsearch":
		h.handleTVSearch(w, r)
	case "get":
		h.handleGet(w, r)
	default:
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><error code="202" description="No such function"/>`))
	}
}
