package download

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// min returns the smaller of a and b.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type FFmpegProgress struct {
	TimeSeconds float64
	SizeBytes   int64
}

var (
	reTime = regexp.MustCompile(`time=(\d+):(\d+):(\d+)\.(\d+)`)
	reSize = regexp.MustCompile(`size=\s*(\d+)kB`)
)

func parseProgress(line string) (FFmpegProgress, bool) {
	var p FFmpegProgress
	tm := reTime.FindStringSubmatch(line)
	sm := reSize.FindStringSubmatch(line)

	if tm == nil || sm == nil {
		return p, false
	}

	h, _ := strconv.ParseFloat(tm[1], 64)
	m, _ := strconv.ParseFloat(tm[2], 64)
	s, _ := strconv.ParseFloat(tm[3], 64)
	ms, _ := strconv.ParseFloat(tm[4], 64)
	p.TimeSeconds = h*3600 + m*60 + s + ms/100

	kb, _ := strconv.ParseInt(sm[1], 10, 64)
	p.SizeBytes = kb * 1024

	return p, true
}

type FFmpegJob struct {
	StreamURL  string
	OutputPath string
	OnProgress func(FFmpegProgress)
}

// resolveHLSVariant fetches the master playlist, finds the highest-bandwidth
// variant, then probes for an unlisted 1080p variant (video=12000000) which
// BBC hosts but omits from the manifest. Falls back to the highest listed
// variant if the 1080p probe returns non-200.
func resolveHLSVariant(masterURL string) string {
	resp, err := http.Get(masterURL)
	if err != nil {
		log.Printf("failed to fetch master playlist: %v", err)
		return masterURL
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read master playlist: %v", err)
		return masterURL
	}
	log.Printf("master playlist fetched: %d bytes, %d lines", len(body), strings.Count(string(body), "\n"))
	log.Printf("master playlist content:\n%s", string(body))

	lines := strings.Split(string(body), "\n")
	bestBW := 0
	bestURL := ""
	bwRe := regexp.MustCompile(`BANDWIDTH=(\d+)`)
	for i, line := range lines {
		if !strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			continue
		}
		if m := bwRe.FindStringSubmatch(line); m != nil {
			bw, _ := strconv.Atoi(m[1])
			if bw > bestBW && i+1 < len(lines) {
				bestBW = bw
				bestURL = strings.TrimSpace(lines[i+1])
			}
		}
	}

	log.Printf("best variant: bw=%d url=%q", bestBW, bestURL)
	if bestURL == "" {
		log.Printf("no variant found in master playlist, returning master URL")
		return masterURL
	}

	// Resolve relative to master playlist base
	base := masterURL
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[:idx+1]
	}
	if !strings.HasPrefix(bestURL, "http") {
		bestURL = base + bestURL
	}

	// BBC hosts unlisted 1080p variants at video=12000000 that aren't in
	// the manifest. Probe for it and use if available.
	if strings.Contains(bestURL, "video=") {
		fhdURL := regexp.MustCompile(`video=\d+`).ReplaceAllString(bestURL, "video=12000000")
		log.Printf("probing unlisted 1080p: %s", fhdURL[:min(len(fhdURL), 120)])
		probeResp, err := http.Head(fhdURL)
		if err != nil {
			log.Printf("1080p probe error: %v", err)
		} else {
			probeResp.Body.Close()
			log.Printf("1080p probe response: %d", probeResp.StatusCode)
			if probeResp.StatusCode == 200 {
				log.Printf("HLS 1080p variant found (unlisted)")
				return fhdURL
			}
		}
	} else {
		log.Printf("best URL does not contain video= segment, skipping 1080p probe")
	}

	log.Printf("HLS variant selected: bandwidth=%d url=%s", bestBW, bestURL[:min(len(bestURL), 100)])
	return bestURL
}

func RunFFmpeg(ctx context.Context, job FFmpegJob) error {
	streamURL := job.StreamURL
	// For HLS master playlists, resolve the highest-bandwidth variant
	// and probe for unlisted 1080p. DASH manifests are handled by ffmpeg.
	if strings.Contains(streamURL, ".m3u8") {
		log.Printf("resolving HLS variant for: %s", streamURL[:min(len(streamURL), 80)])
		streamURL = resolveHLSVariant(streamURL)
		log.Printf("resolved stream URL: %s", streamURL[:min(len(streamURL), 80)])
	} else {
		log.Printf("not HLS, skipping variant resolution: %s", streamURL[:min(len(streamURL), 80)])
	}
	args := []string{
		"-loglevel", "fatal",
		"-stats",
		"-y",
		"-i", streamURL,
		"-c:v", "copy",
		"-c:a", "copy",
		"-bsf:a", "aac_adtstoasc",
		"-movflags", "faststart",
		job.OutputPath,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start ffmpeg: %w", err)
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(scanFFmpegLines)
	for scanner.Scan() {
		line := scanner.Text()
		if prog, ok := parseProgress(line); ok && job.OnProgress != nil {
			job.OnProgress(prog)
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return fmt.Errorf("reading ffmpeg stderr: %w", scanErr)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg: %w", err)
	}
	return nil
}

func scanFFmpegLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' || data[i] == '\r' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func CheckFFmpeg() (string, error) {
	out, err := exec.Command("ffmpeg", "-version").Output()
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found: %w", err)
	}
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 3 {
			return parts[2], nil
		}
		return strings.TrimSpace(lines[0]), nil
	}
	return "unknown", nil
}
