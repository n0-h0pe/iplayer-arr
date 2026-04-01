package bbc

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIBLSearch(t *testing.T) {
	fixture, err := os.ReadFile("testdata/ibl_search.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		if q == "" {
			t.Error("missing q param")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer srv.Close()

	ibl := NewIBL(NewClient())
	ibl.BaseURL = srv.URL

	results, err := ibl.Search("doctor who", 1)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len = %d, want 1", len(results))
	}

	r := results[0]
	if r.PID != "b039d07m" {
		t.Errorf("PID = %q", r.PID)
	}
	if r.Title != "The Unquiet Dead" {
		t.Errorf("Title = %q", r.Title)
	}
	if r.Series != 1 {
		t.Errorf("Series = %d", r.Series)
	}
	if r.EpisodeNum != 3 {
		t.Errorf("EpisodeNum = %d", r.EpisodeNum)
	}
	if r.Channel != "BBC One" {
		t.Errorf("Channel = %q", r.Channel)
	}
}
