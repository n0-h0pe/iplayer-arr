package bbc

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMediaSelectorResolve(t *testing.T) {
	fixture, err := os.ReadFile("testdata/mediaselector.xml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write(fixture)
	}))
	defer srv.Close()

	ms := NewMediaSelector(NewClient())
	ms.BaseURL = srv.URL

	streams, err := ms.Resolve("b039d080")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	if len(streams.Video) < 2 {
		t.Fatalf("expected >= 2 video streams, got %d", len(streams.Video))
	}

	// should be sorted best first
	best := streams.Video[0]
	if best.Height != 720 {
		t.Errorf("best height = %d, want 720", best.Height)
	}
	if best.URL == "" {
		t.Error("best URL is empty")
	}

	if streams.SubtitleURL == "" {
		t.Error("expected subtitle URL")
	}
}

func TestMediaSelectorGeoBlocked(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<mediaSelection><error id="geolocation"/></mediaSelection>`))
	}))
	defer srv.Close()

	ms := NewMediaSelector(NewClient())
	ms.BaseURL = srv.URL

	_, err := ms.Resolve("b039d080")
	if err == nil {
		t.Fatal("expected geo-block error")
	}
	if !IsGeoBlocked(err) {
		t.Errorf("expected geo-block error, got: %v", err)
	}
}
