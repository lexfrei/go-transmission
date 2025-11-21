package transmission

import (
	"bytes"
	"io"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/ymz-ncnk/mok"
)

// RoundTripperMock is a mock for http.RoundTripper.
type RoundTripperMock struct {
	mock *mok.Mock
}

// NewRoundTripperMock creates a new RoundTripperMock.
func NewRoundTripperMock() *RoundTripperMock {
	return &RoundTripperMock{mock: mok.New("RoundTripperMock")}
}

// RegisterRoundTrip registers a function for RoundTrip method.
func (m *RoundTripperMock) RegisterRoundTrip(fn func(*http.Request) (*http.Response, error)) *RoundTripperMock {
	m.mock.Register("RoundTrip", fn)
	return m
}

// RoundTrip implements http.RoundTripper.
func (m *RoundTripperMock) RoundTrip(req *http.Request) (*http.Response, error) {
	results, err := m.mock.Call("RoundTrip", req)
	if err != nil {
		return nil, errors.Wrap(err, "mock call failed")
	}

	resp, _ := results[0].(*http.Response)
	respErr, _ := results[1].(error)

	return resp, respErr
}

// CheckCalls verifies all registered calls were made.
func (m *RoundTripperMock) CheckCalls() []mok.MethodCallsInfo {
	return m.mock.CheckCalls()
}

// jsonResponse creates an HTTP response with JSON body.
func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type":              []string{"application/json"},
			"X-Transmission-Session-Id": []string{"test-session-id"},
		},
		Body: io.NopCloser(bytes.NewBufferString(body)),
	}
}

// sessionIDResponse creates a 409 response for session ID.
func sessionIDResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusConflict,
		Header: http.Header{
			"X-Transmission-Session-Id": []string{"test-session-id"},
		},
		Body: io.NopCloser(bytes.NewBufferString("")),
	}
}
