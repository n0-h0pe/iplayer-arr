package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOpenClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")

	s, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if s == nil {
		t.Fatal("store is nil")
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("db file not created: %v", err)
	}
}

func testStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestDownloadCRUD(t *testing.T) {
	s := testStore(t)

	dl := &Download{
		ID:       "iparr_test1",
		PID:      "b039d07m",
		Title:    "Doctor.Who.S01E03.720p.WEB-DL.AAC.H264-iParr",
		Status:   StatusPending,
		Quality:  "720p",
		Category: "sonarr",
	}

	if err := s.PutDownload(dl); err != nil {
		t.Fatalf("PutDownload: %v", err)
	}

	got, err := s.GetDownload("iparr_test1")
	if err != nil {
		t.Fatalf("GetDownload: %v", err)
	}
	if got.PID != "b039d07m" {
		t.Errorf("PID = %q, want %q", got.PID, "b039d07m")
	}
	if got.Status != StatusPending {
		t.Errorf("Status = %q, want %q", got.Status, StatusPending)
	}

	dl.Status = StatusDownloading
	dl.Progress = 42.5
	if err := s.PutDownload(dl); err != nil {
		t.Fatalf("PutDownload update: %v", err)
	}
	got, _ = s.GetDownload("iparr_test1")
	if got.Progress != 42.5 {
		t.Errorf("Progress = %f, want 42.5", got.Progress)
	}

	all, err := s.ListDownloads()
	if err != nil {
		t.Fatalf("ListDownloads: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("ListDownloads len = %d, want 1", len(all))
	}

	if err := s.DeleteDownload("iparr_test1"); err != nil {
		t.Fatalf("DeleteDownload: %v", err)
	}
	got, err = s.GetDownload("iparr_test1")
	if err != nil {
		t.Fatalf("GetDownload after delete: %v", err)
	}
	if got != nil {
		t.Error("expected nil after delete")
	}
}

func TestDownloadFindByPIDQuality(t *testing.T) {
	s := testStore(t)

	dl := &Download{
		ID:      "iparr_dup1",
		PID:     "b039d07m",
		Quality: "720p",
		Status:  StatusDownloading,
	}
	s.PutDownload(dl)

	found, err := s.FindDownloadByPIDQuality("b039d07m", "720p")
	if err != nil {
		t.Fatalf("FindDownloadByPIDQuality: %v", err)
	}
	if found == nil || found.ID != "iparr_dup1" {
		t.Errorf("expected to find existing download, got %v", found)
	}

	found, _ = s.FindDownloadByPIDQuality("b039d07m", "1080p")
	if found != nil {
		t.Error("expected nil for different quality")
	}
}

func TestHistoryCRUD(t *testing.T) {
	s := testStore(t)

	dl := &Download{ID: "iparr_h1", PID: "test123", Status: StatusCompleted}
	s.PutDownload(dl)
	if err := s.MoveToHistory("iparr_h1"); err != nil {
		t.Fatalf("MoveToHistory: %v", err)
	}

	got, _ := s.GetDownload("iparr_h1")
	if got != nil {
		t.Error("download should be gone from downloads bucket")
	}

	hist, err := s.GetHistory("iparr_h1")
	if err != nil || hist == nil {
		t.Fatalf("GetHistory: err=%v, hist=%v", err, hist)
	}
	if hist.PID != "test123" {
		t.Errorf("PID = %q, want test123", hist.PID)
	}

	all, _ := s.ListHistory()
	if len(all) != 1 {
		t.Errorf("ListHistory len = %d, want 1", len(all))
	}

	s.DeleteHistory("iparr_h1")
	hist, _ = s.GetHistory("iparr_h1")
	if hist != nil {
		t.Error("expected nil after delete")
	}
}

func TestProgrammeCache(t *testing.T) {
	s := testStore(t)

	p := &Programme{
		PID:      "b039d07m",
		Name:     "Doctor Who",
		Episode:  "The Unquiet Dead",
		CachedAt: time.Now(),
	}
	s.PutProgramme(p)

	got, _ := s.GetProgramme("b039d07m")
	if got == nil || got.Name != "Doctor Who" {
		t.Errorf("GetProgramme = %v", got)
	}

	s.DeleteProgramme("b039d07m")
	got, _ = s.GetProgramme("b039d07m")
	if got != nil {
		t.Error("expected nil after delete")
	}
}

func TestConfig(t *testing.T) {
	s := testStore(t)

	s.SetConfig("api_key", "test-key-123")
	val, _ := s.GetConfig("api_key")
	if val != "test-key-123" {
		t.Errorf("GetConfig = %q", val)
	}

	val, _ = s.GetConfig("nonexistent")
	if val != "" {
		t.Errorf("expected empty for nonexistent, got %q", val)
	}
}

func TestSeriesMapping(t *testing.T) {
	s := testStore(t)

	m := &SeriesMapping{TVDBId: "78804", ShowName: "Doctor Who"}
	s.PutSeriesMapping(m)

	got, _ := s.GetSeriesMapping("78804")
	if got == nil || got.ShowName != "Doctor Who" {
		t.Errorf("GetSeriesMapping = %v", got)
	}
}

func TestOverrides(t *testing.T) {
	s := testStore(t)

	o := &ShowOverride{
		ShowName:       "blue peter",
		ForceSeriesNum: 1,
	}
	s.PutOverride(o)

	got, _ := s.GetOverride("blue peter")
	if got == nil || got.ForceSeriesNum != 1 {
		t.Errorf("GetOverride = %v", got)
	}

	all, _ := s.ListOverrides()
	if len(all) != 1 {
		t.Errorf("ListOverrides len = %d", len(all))
	}

	s.DeleteOverride("blue peter")
	got, _ = s.GetOverride("blue peter")
	if got != nil {
		t.Error("expected nil after delete")
	}
}
