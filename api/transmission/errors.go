package transmission

import (
	"net/http"

	"github.com/cockroachdb/errors"
)

// Sentinel errors for client state.
var (
	// ErrNilClient is returned when attempting to use a nil client.
	ErrNilClient = errors.New("client is nil")

	// ErrClosed is returned when attempting to use a closed client.
	ErrClosed = errors.New("client is closed")

	// ErrInvalidURL is returned when the provided URL is invalid.
	ErrInvalidURL = errors.New("invalid URL")
)

// Sentinel errors for RPC operations.
var (
	// ErrRPCFailed is returned when the RPC call fails.
	ErrRPCFailed = errors.New("RPC call failed")

	// ErrCSRFMissing is returned when CSRF token is required but not provided.
	ErrCSRFMissing = errors.New("CSRF session ID missing in response")

	// ErrCSRFRetryFailed is returned when retry after CSRF also fails.
	ErrCSRFRetryFailed = errors.New("CSRF retry failed")
)

// Sentinel errors for HTTP operations.
var (
	// ErrHTTPFailed is returned when HTTP request fails.
	ErrHTTPFailed = errors.New("HTTP request failed")

	// ErrUnauthorized is returned on HTTP 401.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned on HTTP 403.
	ErrForbidden = errors.New("forbidden")

	// ErrNotFound is returned on HTTP 404.
	ErrNotFound = errors.New("not found")

	// ErrConflict is returned on HTTP 409 (CSRF).
	ErrConflict = errors.New("conflict")

	// ErrServerError is returned on HTTP 5xx.
	ErrServerError = errors.New("server error")
)

// Sentinel errors for request/response handling.
var (
	// ErrMarshalRequest is returned when request marshaling fails.
	ErrMarshalRequest = errors.New("failed to marshal request")

	// ErrUnmarshalResponse is returned when response unmarshaling fails.
	ErrUnmarshalResponse = errors.New("failed to unmarshal response")

	// ErrCreateRequest is returned when HTTP request creation fails.
	ErrCreateRequest = errors.New("failed to create HTTP request")

	// ErrReadResponse is returned when reading response body fails.
	ErrReadResponse = errors.New("failed to read response body")

	// ErrNoResponse is returned when the server returns an empty response.
	ErrNoResponse = errors.New("empty response from server")
)

// Sentinel errors for torrent operations.
var (
	// ErrTorrentNotFound is returned when torrent is not found.
	ErrTorrentNotFound = errors.New("torrent not found")

	// ErrDuplicateTorrent is returned when adding a torrent that already exists.
	ErrDuplicateTorrent = errors.New("duplicate torrent")

	// ErrInvalidTorrent is returned when torrent data is invalid.
	ErrInvalidTorrent = errors.New("invalid torrent")
)

// RPCError represents an error returned by the Transmission RPC API.
type RPCError struct {
	Result string
}

// Error implements the error interface.
func (e *RPCError) Error() string {
	return "RPC error: " + e.Result
}

// NewRPCError creates a new RPC error with stack trace.
func NewRPCError(result string) error {
	return errors.WithStack(errors.Mark(&RPCError{Result: result}, ErrRPCFailed))
}

// IsRPCError returns true if the error is an RPCError.
func IsRPCError(err error) bool {
	return errors.Is(err, ErrRPCFailed)
}

// GetRPCError extracts RPCError from error chain if present.
func GetRPCError(err error) *RPCError {
	var rpcErr *RPCError
	if errors.As(err, &rpcErr) {
		return rpcErr
	}
	return nil
}

// HTTPError represents an HTTP-level error.
type HTTPError struct {
	StatusCode int
	Status     string
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	return "HTTP error: " + e.Status
}

// NewHTTPError creates a new HTTP error with appropriate sentinel marker.
func NewHTTPError(statusCode int, status string) error {
	httpErr := &HTTPError{
		StatusCode: statusCode,
		Status:     status,
	}

	// Mark with appropriate sentinel error
	var sentinel error
	switch statusCode {
	case http.StatusUnauthorized:
		sentinel = ErrUnauthorized
	case http.StatusForbidden:
		sentinel = ErrForbidden
	case http.StatusNotFound:
		sentinel = ErrNotFound
	case http.StatusConflict:
		sentinel = ErrConflict
	default:
		if statusCode >= 500 {
			sentinel = ErrServerError
		} else {
			sentinel = ErrHTTPFailed
		}
	}

	return errors.WithStack(errors.Mark(httpErr, sentinel))
}

// IsHTTPError returns true if the error is an HTTPError.
func IsHTTPError(err error) bool {
	var httpErr *HTTPError
	return errors.As(err, &httpErr)
}

// GetHTTPError extracts HTTPError from error chain if present.
func GetHTTPError(err error) *HTTPError {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr
	}
	return nil
}

// IsUnauthorized returns true if the error is HTTP 401 Unauthorized.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden returns true if the error is HTTP 403 Forbidden.
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsNotFound returns true if the error is HTTP 404 Not Found.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsServerError returns true if the error is HTTP 5xx.
func IsServerError(err error) bool {
	return errors.Is(err, ErrServerError)
}
