package transmission

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"

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
	}

	return transport, nil
}

// Do performs a JSON-RPC request with automatic CSRF handling.
func (t *httpTransport) Do(ctx context.Context, method string, args any) (json.RawMessage, error) {
	// Build request
	tag := t.tagSeq.Add(1)
	reqBody := rpcRequest{
		Method:    method,
		Arguments: args,
		Tag:       tag,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.Wrap(err, ErrMarshalRequest.Error())
	}

	// First attempt
	resp, err := t.doRequest(ctx, body)
	if err != nil {
		return nil, err
	}

	// Handle CSRF (HTTP 409)
	if resp.StatusCode == http.StatusConflict {
		// Extract session ID from response
		sessionID := resp.Header.Get(sessionIDHeader)
		_ = resp.Body.Close()

		if sessionID == "" {
			return nil, errors.WithStack(ErrCSRFMissing)
		}

		// Store session ID and retry
		t.sessionID.Store(sessionID)
		resp, err = t.doRequest(ctx, body)
		if err != nil {
			return nil, err
		}
	}
	defer func() { _ = resp.Body.Close() }()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return nil, NewHTTPError(resp.StatusCode, resp.Status)
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, ErrReadResponse.Error())
	}

	// Parse response
	var rpcResp rpcResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	// Check RPC result
	if rpcResp.Result != "success" {
		return nil, NewRPCError(rpcResp.Result)
	}

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
