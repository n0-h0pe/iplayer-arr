package download

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

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

func RunFFmpeg(ctx context.Context, job FFmpegJob) error {
	args := []string{
		"-loglevel", "fatal",
		"-stats",
		"-y",
		"-i", job.StreamURL,
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
		return strings.TrimSpace(lines[0]), nil
	}
	return "unknown", nil
}
