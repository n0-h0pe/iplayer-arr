package download

import (
	"path/filepath"
	"testing"

	"github.com/GiteaLN/iplayer-arr/internal/store"
)

func TestManagerEnqueueAndList(t *testing.T) {
	dir := t.TempDir()
	st, _ := store.Open(filepath.Join(dir, "test.db"))
	defer st.Close()

	m := NewManager(st, filepath.Join(dir, "downloads"), 2)

	id, err := m.Enqueue("b039d07m", "720p", "Test.S01E01.720p", "sonarr")
	if err != nil {
		t.Fatalf("Enqueue: %v", err)
	}
	if id == "" {
		t.Fatal("empty id")
	}

	dl, _ := st.GetDownload(id)
	if dl == nil {
		t.Fatal("download not found in store")
	}
	if dl.Status != store.StatusPending {
		t.Errorf("status = %q, want pending", dl.Status)
	}
	if dl.Category != "sonarr" {
		t.Errorf("category = %q", dl.Category)
	}
}

func TestManagerDeduplicate(t *testing.T) {
	dir := t.TempDir()
	st, _ := store.Open(filepath.Join(dir, "test.db"))
	defer st.Close()

	m := NewManager(st, filepath.Join(dir, "downloads"), 2)

	id1, _ := m.Enqueue("b039d07m", "720p", "Test.S01E01.720p", "sonarr")
	id2, _ := m.Enqueue("b039d07m", "720p", "Test.S01E01.720p", "sonarr")

	if id1 != id2 {
		t.Errorf("duplicate enqueue should return same ID: %q != %q", id1, id2)
	}
}
