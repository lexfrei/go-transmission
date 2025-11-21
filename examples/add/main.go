package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/lexfrei/go-transmission/api/transmission"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: add <torrent-file>")
	}

	url := os.Getenv("TRANSMISSION_URL")
	if url == "" {
		url = "http://localhost:9091/transmission/rpc"
	}

	client, err := transmission.New(url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = client.Close() }()

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	metainfo := base64.StdEncoding.EncodeToString(data)
	result, err := client.TorrentAdd(context.Background(), &transmission.TorrentAddArgs{
		Metainfo: &metainfo,
	})
	if err != nil {
		log.Fatal(err)
	}

	if result.TorrentAdded != nil {
		fmt.Printf("Added: %s (ID: %d)\n", result.TorrentAdded.Name, result.TorrentAdded.ID)
	} else if result.TorrentDuplicate != nil {
		fmt.Printf("Duplicate: %s (ID: %d)\n", result.TorrentDuplicate.Name, result.TorrentDuplicate.ID)
	}
}
