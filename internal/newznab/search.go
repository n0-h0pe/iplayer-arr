package newznab

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/GiteaLN/iplayer-arr/internal/bbc"
	"github.com/GiteaLN/iplayer-arr/internal/store"
)

func (h *Handler) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" || h.ibl == nil {
		writeEmptyRSS(w)
		return
	}

	results, err := h.ibl.Search(q, 1)
	if err != nil {
		writeEmptyRSS(w)
		return
	}

	h.writeResultsRSS(w, r, results, 0, 0)
}

func (h *Handler) handleTVSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	tvdbid := r.URL.Query().Get("tvdbid")
	seasonStr := r.URL.Query().Get("season")
	epStr := r.URL.Query().Get("ep")

	if q == "" && tvdbid != "" && h.store != nil {
		mapping, _ := h.store.GetSeriesMapping(tvdbid)
		if mapping != nil {
			q = mapping.ShowName
		}
	}

	if q == "" || h.ibl == nil {
		writeEmptyRSS(w)
		return
	}

	results, err := h.ibl.Search(q, 1)
	if err != nil {
		writeEmptyRSS(w)
		return
	}

	season, _ := strconv.Atoi(seasonStr)
	ep, _ := strconv.Atoi(epStr)

	h.writeResultsRSS(w, r, results, season, ep)
}

func (h *Handler) writeResultsRSS(w http.ResponseWriter, r *http.Request, results []bbc.IBLResult, filterSeason, filterEp int) {
	var items []string

	for _, res := range results {
		prog := iblResultToProgramme(res)

		if filterSeason > 0 && prog.Series != filterSeason {
			continue
		}
		if filterEp > 0 && prog.EpisodeNum != filterEp {
			continue
		}

		var override *store.ShowOverride
		if h.store != nil {
			override, _ = h.store.GetOverride(prog.Name)
		}

		qualities := []string{"720p"}
		if len(prog.Qualities) > 0 {
			qualities = nil
			for _, q := range prog.Qualities {
				qualities = append(qualities, q.Tag)
			}
		}

		for _, qual := range qualities {
			title, tier := GenerateTitle(prog, qual, override)
			guid := EncodeGUID(res.PID, qual, "original")

			cat := "5040"
			if qual == "540p" || qual == "396p" {
				cat = "5030"
			}

			size := estimateSize(prog.Duration, qual)
			prog.IdentityTier = tier

			pubDate := time.Now().Format(time.RFC1123Z)
			if res.AirDate != "" {
				if t, err := time.Parse("2006-01-02", res.AirDate); err == nil {
					pubDate = t.Format(time.RFC1123Z)
				}
			}

			item := fmt.Sprintf(`    <item>
      <title>%s</title>
      <guid isPermaLink="true">%s/newznab/api?t=get&amp;id=%s</guid>
      <link>%s/newznab/api?t=get&amp;id=%s</link>
      <pubDate>%s</pubDate>
      <enclosure url="%s/newznab/api?t=get&amp;id=%s" length="%d" type="application/x-nzb" />
      <newznab:attr name="category" value="%s" />
      <newznab:attr name="size" value="%d" />`,
				html.EscapeString(title), baseURL(r), guid, baseURL(r), guid, pubDate,
				baseURL(r), guid, size, cat, size)

			if tier == store.TierManual {
				item += `
      <newznab:attr name="iparr:manual" value="true" />`
			}

			item += "\n    </item>"
			items = append(items, item)
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprintf(w, `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:newznab="http://www.newznab.com/DTD/2010/feeds/attributes/">
  <channel>
    <title>iplayer-arr</title>
%s
  </channel>
</rss>`, strings.Join(items, "\n"))
}

func writeEmptyRSS(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:newznab="http://www.newznab.com/DTD/2010/feeds/attributes/">
  <channel><title>iplayer-arr</title></channel>
</rss>`))
}

func iblResultToProgramme(r bbc.IBLResult) *store.Programme {
	return &store.Programme{
		PID:        r.PID,
		Name:       r.Title,
		Episode:    r.Subtitle,
		Series:     r.Series,
		EpisodeNum: r.EpisodeNum,
		Position:   r.Position,
		AirDate:    r.AirDate,
		Channel:    r.Channel,
		Thumbnail:  r.Thumbnail,
	}
}

func estimateSize(durationSec int, quality string) int64 {
	if durationSec == 0 {
		durationSec = 3600
	}
	kbps := map[string]int{
		"1080p": 8500,
		"720p":  5000,
		"540p":  2500,
		"396p":  1250,
	}
	rate, ok := kbps[quality]
	if !ok {
		rate = 5000
	}
	return int64(durationSec) * int64(rate) * 1000 / 8
}

func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}
