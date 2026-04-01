package bbc

import "testing"

func TestRandomUA(t *testing.T) {
	ua := RandomUserAgent()
	if ua == "" {
		t.Fatal("empty user agent")
	}

	// should contain browser identifier
	found := false
	for _, sig := range []string{"Chrome", "Firefox", "Safari", "Edge"} {
		if contains(ua, sig) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("UA %q doesn't look like a browser UA", ua)
	}

	// calling twice should eventually return different UAs (probabilistic)
	seen := map[string]bool{}
	for i := 0; i < 50; i++ {
		seen[RandomUserAgent()] = true
	}
	if len(seen) < 2 {
		t.Error("expected at least 2 different UAs from 50 calls")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
