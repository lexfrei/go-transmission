// Package transmission provides a Go client for the Transmission BitTorrent
// client RPC API.
//
// This package implements all 24 methods of Transmission RPC protocol version 18
// (semver 5.4.0) as documented in the official specification.
//
// # Quick Start
//
//	client, err := transmission.New("http://localhost:9091/transmission/rpc")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// Get all torrents
//	torrents, err := client.TorrentGet(ctx, []string{"id", "name", "status"}, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, t := range torrents {
//	    fmt.Printf("Torrent: %s\n", t.Name)
//	}
//
// # Authentication
//
// If Transmission requires authentication, use the WithAuth option:
//
//	client, err := transmission.New(
//	    "http://localhost:9091/transmission/rpc",
//	    transmission.WithAuth("username", "password"),
//	)
//
// # CSRF Protection
//
// The client automatically handles Transmission's CSRF protection mechanism.
// On receiving HTTP 409, the client extracts the session ID from response headers
// and retries the request automatically.
//
// # RPC Version Compatibility
//
// This client is designed for Transmission 4.1.0+ (RPC version 18).
// Most methods are backwards compatible with earlier versions.
// Use SessionGet to check the server's RPC version for compatibility.
//
// # Reference
//
// The complete API specification is available in openrpc.json in this package.
// Official Transmission RPC documentation:
// https://github.com/transmission/transmission/blob/main/docs/rpc-spec.md
package transmission
