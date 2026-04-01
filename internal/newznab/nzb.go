package newznab

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type GrabInfo struct {
	PID     string
	Quality string
	Version string
}

func EncodeGUID(pid, quality, version string) string {
	raw := fmt.Sprintf("%s:%s:%s", pid, quality, version)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func DecodeGUID(guid string) (*GrabInfo, error) {
	data, err := base64.URLEncoding.DecodeString(guid)
	if err != nil {
		return nil, fmt.Errorf("invalid GUID encoding")
	}
	parts := strings.SplitN(string(data), ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid GUID format")
	}
	return &GrabInfo{PID: parts[0], Quality: parts[1], Version: parts[2]}, nil
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	info, err := DecodeGUID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><error code="300" description="Item not found"/>`))
		return
	}

	downloadID := fmt.Sprintf("%s:%s", info.PID, info.Quality)

	nzb := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE nzb PUBLIC "-//newzBin//DTD NZB 1.1//EN" "http://www.newzbin.com/DTD/nzb/nzb-1.1.dtd">
<nzb>
  <head><meta type="name">iParr Internal</meta></head>
  <file subject="iParr download">
    <groups><group>iparr.internal</group></groups>
    <segments><segment number="1">%s</segment></segments>
  </file>
</nzb>`, downloadID)

	w.Header().Set("Content-Type", "application/x-nzb")
	w.Write([]byte(nzb))
}
