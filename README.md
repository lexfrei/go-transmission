# go-transmission

Go client library for [Transmission](https://transmissionbt.com) BitTorrent client RPC API.

## Features

- Full implementation of Transmission RPC protocol version 18 (semver 5.4.0)
- All 24 RPC methods supported
- Automatic CSRF protection handling
- HTTP Basic Authentication support
- Context-based cancellation and timeouts
- Idiomatic Go API with comprehensive types

## Requirements

- Go 1.21 or later
- Transmission 4.0+ (RPC version 17+)

## Installation

```bash
go get github.com/lexfrei/go-transmission/api/transmission
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/lexfrei/go-transmission/api/transmission"
)

func main() {
    // Create client
    client, err := transmission.New("http://localhost:9091/transmission/rpc")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()

    // Get all torrents
    result, err := client.TorrentGet(ctx, []string{"id", "name", "status", "percentDone"}, nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, t := range result.Torrents {
        fmt.Printf("Torrent: %s (%.1f%%)\n", *t.Name, *t.PercentDone*100)
    }
}
```

## Authentication

If Transmission requires authentication:

```go
client, err := transmission.New(
    "http://localhost:9091/transmission/rpc",
    transmission.WithAuth("username", "password"),
)
```

## Timeout

Set a custom timeout:

```go
client, err := transmission.New(
    "http://localhost:9091/transmission/rpc",
    transmission.WithTimeout(30 * time.Second),
)
```

## API Methods

### Torrent Actions
- `TorrentStart` - Start torrents
- `TorrentStartNow` - Start torrents immediately (bypass queue)
- `TorrentStop` - Stop torrents
- `TorrentVerify` - Verify torrent data
- `TorrentReannounce` - Force tracker announce

### Torrent Accessors
- `TorrentGet` - Get torrent information
- `TorrentGetByHash` - Get torrents by hash
- `TorrentGetRecentlyActive` - Get recently active torrents

### Torrent Mutators
- `TorrentSet` - Modify torrent properties
- `TorrentAdd` - Add a new torrent
- `TorrentRemove` - Remove torrents
- `TorrentSetLocation` - Move torrent data
- `TorrentRenamePath` - Rename files in torrent

### Session
- `SessionGet` - Get session configuration
- `SessionSet` - Modify session configuration
- `SessionStats` - Get transfer statistics
- `SessionClose` - Shutdown daemon

### Queue Management
- `QueueMoveTop` - Move to top of queue
- `QueueMoveUp` - Move up in queue
- `QueueMoveDown` - Move down in queue
- `QueueMoveBottom` - Move to bottom of queue

### System
- `BlocklistUpdate` - Update IP blocklist
- `PortTest` - Test port accessibility
- `FreeSpace` - Check available disk space

### Bandwidth Groups
- `GroupSet` - Create/modify bandwidth group
- `GroupGet` - Get bandwidth group info

## OpenRPC Specification

This project includes a complete [OpenRPC 1.3.2](https://spec.open-rpc.org/) specification for the Transmission RPC API:

**[docs/transmission-rpc-v5.4.0.openrpc.json](docs/transmission-rpc-v5.4.0.openrpc.json)**

- **RPC Version**: 18 (semver 5.4.0)
- **Transmission**: 4.1.0+
- **Coverage**: All 24 methods, 23 schemas, with examples

The OpenRPC spec can be used independently from this Go library — for code generation in other languages, API documentation, or tooling integration.

For official Transmission RPC documentation, see:
https://github.com/transmission/transmission/blob/main/docs/rpc-spec.md

## License

MIT License - see [LICENSE](LICENSE) for details.
