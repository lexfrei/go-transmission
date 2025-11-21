package transmission

import (
	"net/http"
	"time"
)

// config holds the client configuration.
type config struct {
	url        string
	username   string
	password   string
	timeout    time.Duration
	httpClient *http.Client
}

// Option is a function that configures the client.
type Option func(*config)

// WithAuth configures HTTP Basic Authentication.
func WithAuth(username, password string) Option {
	return func(c *config) {
		c.username = username
		c.password = password
	}
}

// WithTimeout sets the HTTP client timeout.
// Default is no timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client.
// This overrides the timeout option.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *config) {
		c.httpClient = httpClient
	}
}
