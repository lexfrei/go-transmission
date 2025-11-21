package transmission

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

// sessionGetArgs is used for session-get method.
type sessionGetArgs struct {
	Fields []string `json:"fields,omitempty"`
}

// SessionGet retrieves session configuration.
func (c *client) SessionGet(ctx context.Context, fields []string) (*Session, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	var args any
	if len(fields) > 0 {
		args = sessionGetArgs{Fields: fields}
	}

	resp, err := c.transport.Do(ctx, "session-get", args)
	if err != nil {
		return nil, err
	}

	var result Session
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return &result, nil
}

// SessionSet modifies session configuration.
func (c *client) SessionSet(ctx context.Context, args *SessionSetArgs) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	_, err := c.transport.Do(ctx, "session-set", args)
	return err
}

// SessionStats retrieves session statistics.
func (c *client) SessionStats(ctx context.Context) (*SessionStats, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	resp, err := c.transport.Do(ctx, "session-stats", nil)
	if err != nil {
		return nil, err
	}

	var result SessionStats
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return &result, nil
}

// SessionClose gracefully shuts down the Transmission daemon.
func (c *client) SessionClose(ctx context.Context) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	_, err := c.transport.Do(ctx, "session-close", nil)
	return err
}
