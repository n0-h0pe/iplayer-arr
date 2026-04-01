package download

import "testing"

func TestParseFFmpegProgress(t *testing.T) {
	tests := []struct {
		line     string
		wantTime float64
		wantSize int64
		wantOK   bool
	}{
		{
			"frame=  720 fps= 25 q=-1.0 size=  456789kB time=00:12:34.56 bitrate=5000.0kbits/s speed=2.5x",
			754.56, 456789 * 1024, true,
		},
		{
			"size=  123456kB time=00:05:00.00 bitrate=3300.0kbits/s speed=1.2x",
			300.0, 123456 * 1024, true,
		},
		{
			"some random line",
			0, 0, false,
		},
	}

	for _, tt := range tests {
		prog, ok := parseProgress(tt.line)
		if ok != tt.wantOK {
			t.Errorf("parseProgress(%q): ok = %v, want %v", tt.line, ok, tt.wantOK)
			continue
		}
		if !ok {
			continue
		}
		if prog.TimeSeconds != tt.wantTime {
			t.Errorf("time = %f, want %f", prog.TimeSeconds, tt.wantTime)
		}
		if prog.SizeBytes != tt.wantSize {
			t.Errorf("size = %d, want %d", prog.SizeBytes, tt.wantSize)
		}
	}
}
