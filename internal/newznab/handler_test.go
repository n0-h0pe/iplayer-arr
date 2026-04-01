package newznab

import (
	"encoding/xml"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCapsEndpoint(t *testing.T) {
	h := NewHandler(nil, nil, nil)
	req := httptest.NewRequest("GET", "/newznab/api?t=caps", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("status = %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, `<searching>`) {
		t.Error("missing <searching> in caps")
	}
	if !strings.Contains(body, `supportedParams="q,season,ep,tvdbid"`) {
		t.Error("missing tvsearch supportedParams")
	}
	if !strings.Contains(body, `id="5000"`) {
		t.Error("missing TV category 5000")
	}

	var caps struct{}
	if err := xml.Unmarshal(w.Body.Bytes(), &caps); err != nil {
		t.Errorf("invalid XML: %v", err)
	}
}
