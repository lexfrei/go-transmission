package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/lexfrei/go-transmission/api/transmission"
)

// slogAdapter adapts slog.Logger to transmission.Logger interface.
type slogAdapter struct {
	logger *slog.Logger
}

func (s *slogAdapter) Debug(msg string, fields ...transmission.Field) {
	s.logger.Debug(msg, fieldsToAttrs(fields)...)
}

func (s *slogAdapter) Warn(msg string, fields ...transmission.Field) {
	s.logger.Warn(msg, fieldsToAttrs(fields)...)
}

func (s *slogAdapter) Error(msg string, fields ...transmission.Field) {
	s.logger.Error(msg, fieldsToAttrs(fields)...)
}

func fieldsToAttrs(fields []transmission.Field) []any {
	attrs := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		attrs = append(attrs, f.Key, f.Value)
	}

	return attrs
}

func main() {
	// Get URL from environment or use default
	url := os.Getenv("TRANSMISSION_URL")
	if url == "" {
		url = "http://localhost:9091/transmission/rpc"
	}

	// Create slog logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create client with optional authentication
	var opts []transmission.Option
	if user := os.Getenv("TRANSMISSION_USER"); user != "" {
		pass := os.Getenv("TRANSMISSION_PASS")
		opts = append(opts, transmission.WithAuth(user, pass))
	}
	opts = append(opts, transmission.WithTimeout(30*time.Second))
	opts = append(opts, transmission.WithLogger(&slogAdapter{logger: logger}))

	client, err := transmission.New(url, opts...)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()

	// Get session info
	session, err := client.SessionGet(ctx, []string{"version", "rpc-version", "download-dir"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get session: %v\n", err)
		return
	}

	fmt.Printf("Transmission version: %s (RPC v%d)\n", *session.Version, *session.RPCVersion)
	fmt.Printf("Download directory: %s\n", *session.DownloadDir)
	fmt.Println()

	// Get session stats
	stats, err := client.SessionStats(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get stats: %v\n", err)
		return
	}

	fmt.Printf("Active torrents: %d\n", stats.ActiveTorrentCount)
	fmt.Printf("Total torrents: %d\n", stats.TorrentCount)
	fmt.Printf("Download speed: %.2f KB/s\n", float64(stats.DownloadSpeed)/1024)
	fmt.Printf("Upload speed: %.2f KB/s\n", float64(stats.UploadSpeed)/1024)
	fmt.Println()

	// Get all torrents
	result, err := client.TorrentGet(ctx, []string{
		"id", "name", "status", "percentDone",
		"rateDownload", "rateUpload", "eta",
	}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get torrents: %v\n", err)
		return
	}

	if len(result.Torrents) == 0 {
		fmt.Println("No torrents found.")
		return
	}

	fmt.Printf("Torrents (%d):\n", len(result.Torrents))
	for i := range result.Torrents {
		t := &result.Torrents[i]
		status := "Unknown"
		if t.Status != nil {
			status = t.Status.String()
		}

		name := "Unknown"
		if t.Name != nil {
			name = *t.Name
		}

		percent := 0.0
		if t.PercentDone != nil {
			percent = *t.PercentDone * 100
		}

		fmt.Printf("  [%d] %s - %s (%.1f%%)\n", *t.ID, name, status, percent)

		// Show speeds for active torrents
		if t.RateDownload != nil && *t.RateDownload > 0 {
			fmt.Printf("       Download: %.2f KB/s\n", float64(*t.RateDownload)/1024)
		}
		if t.RateUpload != nil && *t.RateUpload > 0 {
			fmt.Printf("       Upload: %.2f KB/s\n", float64(*t.RateUpload)/1024)
		}
	}
}
