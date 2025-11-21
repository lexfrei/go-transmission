package transmission

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

// blocklistUpdateResult is the response from blocklist-update.
type blocklistUpdateResult struct {
	BlocklistSize int `json:"blocklist-size"`
}

// BlocklistUpdate updates the IP blocklist.
func (c *client) BlocklistUpdate(ctx context.Context) (int, error) {
	if err := c.checkClosed(); err != nil {
		return 0, err
	}

	resp, err := c.transport.Do(ctx, "blocklist-update", nil)
	if err != nil {
		return 0, err
	}

	var result blocklistUpdateResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return 0, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return result.BlocklistSize, nil
}

// portTestResult is the response from port-test.
type portTestResult struct {
	PortIsOpen bool `json:"port-is-open"`
}

// PortTest tests if the peer port is accessible.
func (c *client) PortTest(ctx context.Context) (bool, error) {
	if err := c.checkClosed(); err != nil {
		return false, err
	}

	resp, err := c.transport.Do(ctx, "port-test", nil)
	if err != nil {
		return false, err
	}

	var result portTestResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return false, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return result.PortIsOpen, nil
}

// freeSpaceArgs is used for free-space method.
type freeSpaceArgs struct {
	Path string `json:"path"`
}

// FreeSpace checks available disk space at the specified path.
func (c *client) FreeSpace(ctx context.Context, path string) (*FreeSpace, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	args := freeSpaceArgs{Path: path}

	resp, err := c.transport.Do(ctx, "free-space", args)
	if err != nil {
		return nil, err
	}

	var result FreeSpace
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return &result, nil
}
