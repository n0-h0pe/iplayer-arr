package bbc

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ttDoc struct {
	XMLName xml.Name `xml:"tt"`
	Body    struct {
		Divs []struct {
			Paragraphs []struct {
				Begin string `xml:"begin,attr"`
				End   string `xml:"end,attr"`
				Text  string `xml:",chardata"`
				Spans []struct {
					Text string `xml:",chardata"`
				} `xml:"span"`
			} `xml:"p"`
		} `xml:"div"`
	} `xml:"body"`
}

// TTMLToSRT converts BBC TTML subtitle data to SRT format.
// BBC TTML may use period-based timing (HH:MM:SS.mmm) or frame-based timing
// (HH:MM:SS:FF at 25fps). Text may be wrapped in <span> elements.
func TTMLToSRT(ttml []byte) ([]byte, error) {
	var doc ttDoc
	if err := xml.Unmarshal(ttml, &doc); err != nil {
		return nil, fmt.Errorf("parse TTML: %w", err)
	}

	var b strings.Builder
	index := 1

	for _, div := range doc.Body.Divs {
		for _, p := range div.Paragraphs {
			text := p.Text
			if text == "" {
				var parts []string
				for _, s := range p.Spans {
					if t := strings.TrimSpace(s.Text); t != "" {
						parts = append(parts, t)
					}
				}
				text = strings.Join(parts, " ")
			}
			text = strings.TrimSpace(stripTags(text))
			if text == "" {
				continue
			}

			begin := toSRTTime(p.Begin)
			end := toSRTTime(p.End)

			fmt.Fprintf(&b, "%d\n%s --> %s\n%s\n\n", index, begin, end, text)
			index++
		}
	}

	return []byte(b.String()), nil
}

var reHTMLTag = regexp.MustCompile(`<[^>]+>`)

func stripTags(s string) string {
	return reHTMLTag.ReplaceAllString(s, "")
}

// toSRTTime converts a TTML timestamp to SRT format.
// TTML uses HH:MM:SS.mmm (period-based) or HH:MM:SS:FF (frame-based at 25fps).
// SRT uses HH:MM:SS,mmm.
func toSRTTime(t string) string {
	t = strings.TrimSpace(t)
	if strings.Count(t, ":") == 3 {
		// frame-based: HH:MM:SS:FF -> convert FF to milliseconds (assume 25fps)
		parts := strings.Split(t, ":")
		if len(parts) == 4 {
			frames, _ := strconv.Atoi(parts[3])
			ms := (frames * 1000) / 25
			return fmt.Sprintf("%s:%s:%s,%03d", parts[0], parts[1], parts[2], ms)
		}
	}
	// period-based: HH:MM:SS.mmm -> HH:MM:SS,mmm
	return strings.Replace(t, ".", ",", 1)
}
