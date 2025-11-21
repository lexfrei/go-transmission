package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/errors"
)

const (
	// sessionIDHeader is the CSRF protection header name.
	sessionIDHeader = "X-Transmission-Session-Id"

	// contentType is the content type for JSON-RPC requests.
	contentType = "application/json"
)

// rpcRequest represents a JSON-RPC request to Transmission.
type rpcRequest struct {
	Method    string `json:"method"`
	Arguments any    `json:"arguments,omitempty"`
	Tag       int64  `json:"tag,omitempty"`
}

// rpcResponse represents a JSON-RPC response from Transmission.
type rpcResponse struct {
	Result    string          `json:"result"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
	Tag       int64           `json:"tag,omitempty"`
}

// httpTransport handles HTTP communication with the Transmission RPC server.
type httpTransport struct {
	client    *http.Client
	url       string
	username  string
	password  string
	sessionID atomic.Value // string
	tagSeq    atomic.Int64
	logger    Logger
}

// newHTTPTransport creates a new HTTP transport.
func newHTTPTransport(cfg *config) (*httpTransport, error) {
	if _, err := url.Parse(cfg.url); err != nil {
		return nil, errors.Wrap(err, ErrInvalidURL.Error())
	}

	httpClient := cfg.httpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: cfg.timeout,
		}
	}

	transport := &httpTransport{
		client:   httpClient,
		url:      cfg.url,
		username: cfg.username,
		password: cfg.password,
		logger:   cfg.logger,
	}

	return transport, nil
}

// Do performs a JSON-RPC request with automatic CSRF handling.
func (t *httpTransport) Do(ctx context.Context, method string, args any) (json.RawMessage, error) {
	tag := t.tagSeq.Add(1)

	body, err := t.marshalRequest(method, args, tag)
	if err != nil {
		return nil, err
	}

	t.logger.Debug("sending RPC request", Field{Key: "method", Value: method}, Field{Key: "tag", Value: tag})

	start := time.Now()

	resp, err := t.doRequestWithCSRF(ctx, method, body)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return t.handleResponse(resp, method, start)
}

// marshalRequest builds and marshals an RPC request.
func (t *httpTransport) marshalRequest(method string, args any, tag int64) ([]byte, error) {
	reqBody := rpcRequest{Method: method, Arguments: args, Tag: tag}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.logger.Error("failed to marshal request",
			Field{Key: "method", Value: method}, Field{Key: "error", Value: err.Error()})

		return nil, errors.Wrap(err, ErrMarshalRequest.Error())
	}

	return body, nil
}

// doRequestWithCSRF performs request with automatic CSRF retry.
func (t *httpTransport) doRequestWithCSRF(ctx context.Context, method string, body []byte) (*http.Response, error) {
	resp, err := t.doRequest(ctx, body)
	if err != nil {
		t.logger.Error("HTTP request failed", Field{Key: "method", Value: method}, Field{Key: "error", Value: err.Error()})

		return nil, err
	}

	if resp.StatusCode != http.StatusConflict {
		return resp, nil
	}

	// Handle CSRF (HTTP 409)
	sessionID := resp.Header.Get(sessionIDHeader)
	_ = resp.Body.Close()

	if sessionID == "" {
		t.logger.Error("CSRF session ID missing in 409 response", Field{Key: "method", Value: method})

		return nil, errors.WithStack(ErrCSRFMissing)
	}

	t.logger.Debug("received new CSRF session ID, retrying request", Field{Key: "method", Value: method})
	t.sessionID.Store(sessionID)

	resp, err = t.doRequest(ctx, body)
	if err != nil {
		t.logger.Error("HTTP request failed after CSRF retry",
			Field{Key: "method", Value: method}, Field{Key: "error", Value: err.Error()})

		return nil, err
	}

	return resp, nil
}

// handleResponse processes the HTTP response and extracts RPC result.
func (t *httpTransport) handleResponse(resp *http.Response, method string, start time.Time) (json.RawMessage, error) {
	duration := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		t.logger.Error("unexpected HTTP status",
			Field{Key: "method", Value: method}, Field{Key: "status", Value: resp.StatusCode},
			Field{Key: "duration_ms", Value: duration.Milliseconds()})

		return nil, NewHTTPError(resp.StatusCode, resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.logger.Error("failed to read response body",
			Field{Key: "method", Value: method}, Field{Key: "error", Value: err.Error()})

		return nil, errors.Wrap(err, ErrReadResponse.Error())
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		t.logger.Error("failed to unmarshal response",
			Field{Key: "method", Value: method}, Field{Key: "error", Value: err.Error()})

		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	if rpcResp.Result != "success" {
		t.logger.Error("RPC error", Field{Key: "method", Value: method},
			Field{Key: "result", Value: rpcResp.Result}, Field{Key: "duration_ms", Value: duration.Milliseconds()})

		return nil, NewRPCError(rpcResp.Result)
	}

	t.logger.Debug("RPC request completed",
		Field{Key: "method", Value: method}, Field{Key: "duration_ms", Value: duration.Milliseconds()})

	return rpcResp.Arguments, nil
}

// doRequest performs the actual HTTP request.
func (t *httpTransport) doRequest(ctx context.Context, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, t.url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, ErrCreateRequest.Error())
	}

	// Set headers
	req.Header.Set("Content-Type", contentType)

	// Add session ID if available
	if sessionID, ok := t.sessionID.Load().(string); ok && sessionID != "" {
		req.Header.Set(sessionIDHeader, sessionID)
	}

	// Add basic auth if configured
	if t.username != "" {
		req.SetBasicAuth(t.username, t.password)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, ErrHTTPFailed.Error())
	}

	return resp, nil
}
