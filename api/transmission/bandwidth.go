package transmission

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

// GroupSet creates or modifies a bandwidth group.
func (c *client) GroupSet(ctx context.Context, args *BandwidthGroup) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	_, err := c.transport.Do(ctx, "group-set", args)
	return err
}

// groupGetArgs is used for group-get method.
type groupGetArgs struct {
	Group any `json:"group,omitempty"`
}

// groupGetResult is the response from group-get.
type groupGetResult struct {
	Group []BandwidthGroup `json:"group"`
}

// GroupGet retrieves bandwidth group information.
func (c *client) GroupGet(ctx context.Context, names []string) ([]BandwidthGroup, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	var args any
	if len(names) > 0 {
		if len(names) == 1 {
			args = groupGetArgs{Group: names[0]}
		} else {
			args = groupGetArgs{Group: names}
		}
	}

	resp, err := c.transport.Do(ctx, "group-get", args)
	if err != nil {
		return nil, err
	}

	var result groupGetResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return result.Group, nil
}
