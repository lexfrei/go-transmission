# OpenRPC Specification for Transmission RPC API

This directory contains a complete OpenRPC 1.3.2 specification for the Transmission BitTorrent client RPC API.

## Overview

**File**: `transmission-rpc-v5.4.0.openrpc.json`

**Coverage**:
- 24 RPC methods covering all Transmission operations
- 23 data type schemas with full type definitions
- Complete examples for every method
- CSRF protection and authentication documentation
- Transmission RPC version 18 (semver 5.4.0)

## Methods Included

### Torrent Actions (5 methods)
- `torrent-start` - Start torrents
- `torrent-start-now` - Start torrents immediately, bypassing queue
- `torrent-stop` - Stop torrents
- `torrent-verify` - Verify torrent data integrity
- `torrent-reannounce` - Force tracker announce

### Torrent Management (4 methods)
- `torrent-set` - Modify torrent properties (20+ parameters)
- `torrent-get` - Retrieve torrent information (70+ available fields)
- `torrent-add` - Add new torrent from file/URL/metainfo
- `torrent-remove` - Remove torrents with optional data deletion

### Torrent Operations (2 methods)
- `torrent-set-location` - Move torrent data to new location
- `torrent-rename-path` - Rename files/directories within torrent

### Session Configuration (3 methods)
- `session-set` - Modify session settings (40+ parameters)
- `session-get` - Retrieve session configuration
- `session-stats` - Get transfer statistics

### Queue Management (4 methods)
- `queue-move-top` - Move torrents to top of queue
- `queue-move-up` - Move torrents up one position
- `queue-move-down` - Move torrents down one position
- `queue-move-bottom` - Move torrents to bottom of queue

### System Utilities (4 methods)
- `blocklist-update` - Update IP blocklist
- `port-test` - Test peer port accessibility
- `session-close` - Gracefully shutdown daemon
- `free-space` - Check available disk space

### Bandwidth Groups (2 methods)
- `group-set` - Create/modify bandwidth group
- `group-get` - Retrieve bandwidth group information

## Code Generation for Go

### Using OpenRPC Generator

The specification can be used with OpenRPC code generators:

```bash
# Install OpenRPC Generator
npm install -g @open-rpc/generator

# Generate Go client
openrpc-generator generate \
  --schema transmission-rpc-v5.4.0.openrpc.json \
  --language go \
  --output-dir ./generated
```

### Alternative: Manual Implementation

You can also use this specification as a reference for manual client implementation:

1. **Type Definitions**: Use `components/schemas` section for Go struct definitions
2. **Method Signatures**: Use `methods` section for API method implementations
3. **Examples**: Reference provided examples for testing

### Key Implementation Notes

#### CSRF Protection
Transmission requires `X-Transmission-Session-Id` header:
```go
// On HTTP 409, extract session ID from response headers
// and retry request with the header
sessionID := resp.Header.Get("X-Transmission-Session-Id")
req.Header.Set("X-Transmission-Session-Id", sessionID)
```

#### Authentication
Optional HTTP Basic Auth:
```go
req.SetBasicAuth(username, password)
```

#### Request Format
All requests use JSON-RPC 2.0 structure:
```json
{
  "method": "torrent-get",
  "arguments": {
    "fields": ["id", "name", "status"]
  },
  "tag": 1
}
```

#### Response Format
```json
{
  "result": "success",
  "arguments": {
    "torrents": [...]
  },
  "tag": 1
}
```

## Validation

The specification is validated:
- ✅ Valid JSON syntax
- ✅ OpenRPC 1.3.2 structure
- ✅ All 24 methods documented
- ✅ All 23 schemas defined
- ✅ Examples provided for all methods

## API Versioning

This specification covers:
- **RPC Version**: 18 (integer)
- **RPC Semver**: 5.4.0
- **Minimum Supported**: RPC version 1
- **Transmission Version**: 4.1.0+

Clients should check `rpc-version` via `session-get` for compatibility.

## Special Features

### Recently Active Torrents
Use `ids: "recently-active"` with `torrent-get` to get:
- Torrents active since last call
- List of removed torrent IDs in `removed` field

### Torrent Identifiers
Methods accept flexible torrent IDs:
- Numeric IDs: `[1, 2, 3]`
- Hash strings: `["abc123def456"]`
- Special value: `"recently-active"`

### Duplicate Handling
`torrent-add` returns:
- `torrent-added` for new torrents
- `torrent-duplicate` for existing torrents
- Both cases return `result: "success"`

## External References

- **Official Specification**: https://github.com/transmission/transmission/blob/main/docs/rpc-spec.md
- **OpenRPC Specification**: https://spec.open-rpc.org/
- **Transmission Project**: https://transmissionbt.com

## License

This OpenRPC specification is provided for the Transmission RPC API, which is licensed under MIT.
The specification itself follows the same licensing terms as the Transmission project.
