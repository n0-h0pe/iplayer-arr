package bbc

import "testing"

func TestTTMLToSRT(t *testing.T) {
	ttml := `<?xml version="1.0" encoding="utf-8"?>
<tt xmlns="http://www.w3.org/ns/ttml" xmlns:ttp="http://www.w3.org/ns/ttml#parameter" ttp:frameRate="25">
  <body>
    <div>
      <p begin="00:00:05.000" end="00:00:08.000">Hello, world!</p>
      <p begin="00:00:10.000" end="00:00:13.500">Second subtitle line.</p>
    </div>
  </body>
</tt>`

	srt, err := TTMLToSRT([]byte(ttml))
	if err != nil {
		t.Fatalf("TTMLToSRT: %v", err)
	}

	expected := "1\n00:00:05,000 --> 00:00:08,000\nHello, world!\n\n2\n00:00:10,000 --> 00:00:13,500\nSecond subtitle line.\n\n"
	if string(srt) != expected {
		t.Errorf("SRT output:\n%s\nwant:\n%s", srt, expected)
	}
}

func TestTTMLToSRTEmpty(t *testing.T) {
	ttml := `<?xml version="1.0"?><tt xmlns="http://www.w3.org/ns/ttml"><body><div></div></body></tt>`
	srt, err := TTMLToSRT([]byte(ttml))
	if err != nil {
		t.Fatalf("TTMLToSRT: %v", err)
	}
	if len(srt) != 0 {
		t.Errorf("expected empty SRT, got %q", srt)
	}
}
