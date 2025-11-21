//go:build ignore

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	rpcURL      = "http://localhost:9091/transmission/rpc"
	responsesDir = "api/transmission/testdata/responses"
)

var sessionID string

func main() {
	if err := os.MkdirAll(responsesDir, 0o755); err != nil {
		panic(err)
	}

	// Get session ID first
	getSessionID()

	// Success responses
	captureResponse("session-get", map[string]any{
		"method": "session-get",
	})

	captureResponse("session-get-fields", map[string]any{
		"method":    "session-get",
		"arguments": map[string]any{"fields": []string{"version", "rpc-version"}},
	})

	captureResponse("session-stats", map[string]any{
		"method": "session-stats",
	})

	captureResponse("free-space", map[string]any{
		"method":    "free-space",
		"arguments": map[string]any{"path": "/config"},
	})

	captureResponse("port-test", map[string]any{
		"method": "port-test",
	})

	captureResponse("group-get", map[string]any{
		"method": "group-get",
	})

	// Add a torrent first
	torrentData, _ := os.ReadFile("api/transmission/testdata/ubuntu-24.04.3-live-server-amd64.iso.torrent")
	metainfo := base64.StdEncoding.EncodeToString(torrentData)

	captureResponse("torrent-add", map[string]any{
		"method": "torrent-add",
		"arguments": map[string]any{
			"metainfo": metainfo,
			"paused":   true,
		},
	})

	// Add duplicate
	captureResponse("torrent-add-duplicate", map[string]any{
		"method": "torrent-add",
		"arguments": map[string]any{
			"metainfo": metainfo,
		},
	})

	captureResponse("torrent-get-all", map[string]any{
		"method": "torrent-get",
		"arguments": map[string]any{
			"fields": []string{"id", "name", "status", "hashString"},
		},
	})

	captureResponse("torrent-get-by-id", map[string]any{
		"method": "torrent-get",
		"arguments": map[string]any{
			"ids":    []int{1},
			"fields": []string{"id", "name", "downloadLimit", "downloadLimited"},
		},
	})

	captureResponse("torrent-get-recently-active", map[string]any{
		"method": "torrent-get",
		"arguments": map[string]any{
			"ids":    "recently-active",
			"fields": []string{"id", "name"},
		},
	})

	// Error responses
	captureResponse("error-free-space-invalid-path", map[string]any{
		"method":    "free-space",
		"arguments": map[string]any{"path": "/nonexistent/path"},
	})

	captureResponse("error-blocklist-update", map[string]any{
		"method": "blocklist-update",
	})

	captureResponse("error-torrent-add-invalid", map[string]any{
		"method": "torrent-add",
		"arguments": map[string]any{
			"metainfo": base64.StdEncoding.EncodeToString([]byte("not a torrent file")),
		},
	})

	captureResponse("error-torrent-add-no-args", map[string]any{
		"method":    "torrent-add",
		"arguments": map[string]any{},
	})

	// Non-existent torrent operations
	captureResponse("torrent-get-nonexistent", map[string]any{
		"method": "torrent-get",
		"arguments": map[string]any{
			"ids":    []int{99999},
			"fields": []string{"id", "name"},
		},
	})

	captureResponse("torrent-start-nonexistent", map[string]any{
		"method": "torrent-start",
		"arguments": map[string]any{
			"ids": []int{99999},
		},
	})

	captureResponse("torrent-stop-nonexistent", map[string]any{
		"method": "torrent-stop",
		"arguments": map[string]any{
			"ids": []int{99999},
		},
	})

	captureResponse("torrent-remove-nonexistent", map[string]any{
		"method": "torrent-remove",
		"arguments": map[string]any{
			"ids": []int{99999},
		},
	})

	captureResponse("torrent-set-nonexistent", map[string]any{
		"method": "torrent-set",
		"arguments": map[string]any{
			"ids":           []int{99999},
			"downloadLimit": 100,
		},
	})

	// Invalid method
	captureResponse("error-invalid-method", map[string]any{
		"method": "not-a-real-method",
	})

	// Clean up - remove the test torrent
	doRequest(map[string]any{
		"method": "torrent-remove",
		"arguments": map[string]any{
			"ids":               []int{1},
			"delete-local-data": true,
		},
	})

	fmt.Println("Done capturing responses!")
}

func getSessionID() {
	req, _ := http.NewRequest("POST", rpcURL, bytes.NewReader([]byte(`{"method":"session-get"}`)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 409 {
		sessionID = resp.Header.Get("X-Transmission-Session-Id")
		fmt.Printf("Got session ID: %s\n", sessionID)
	}
}

func captureResponse(name string, request map[string]any) {
	fmt.Printf("Capturing %s...\n", name)

	body := doRequest(request)

	// Pretty print JSON
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		fmt.Printf("  Warning: could not prettify JSON: %v\n", err)
		prettyJSON.Write(body)
	}

	filename := filepath.Join(responsesDir, name+".json")
	if err := os.WriteFile(filename, prettyJSON.Bytes(), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("  Saved to %s\n", filename)
}

func doRequest(request map[string]any) []byte {
	reqBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", rpcURL, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transmission-Session-Id", sessionID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}
