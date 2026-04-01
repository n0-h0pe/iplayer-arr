package bbc

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestResolveVPID(t *testing.T) {
	fixture, err := os.ReadFile("testdata/playlist.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer srv.Close()

	resolver := NewPlaylistResolver(NewClient())
	resolver.BaseURL = srv.URL

	info, err := resolver.Resolve("b039d07m")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if info.VPID != "b039d080" {
		t.Errorf("VPID = %q, want b039d080", info.VPID)
	}
	if info.Duration != 2700 {
		t.Errorf("Duration = %d, want 2700", info.Duration)
	}
	if len(info.Versions) != 2 {
		t.Errorf("Versions len = %d, want 2", len(info.Versions))
	}
}
