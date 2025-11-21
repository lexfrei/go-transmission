package transmission

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
)

// torrentActionArgs is used for torrent action methods (start, stop, etc.).
type torrentActionArgs struct {
	IDs any `json:"ids,omitempty"`
}

// TorrentStart starts one or more torrents.
func (c *client) TorrentStart(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := torrentActionArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "torrent-start", args)
	return err
}

// TorrentStartNow starts torrents immediately, bypassing the queue.
func (c *client) TorrentStartNow(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := torrentActionArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "torrent-start-now", args)
	return err
}

// TorrentStop stops one or more torrents.
func (c *client) TorrentStop(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := torrentActionArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "torrent-stop", args)
	return err
}

// TorrentVerify verifies local data integrity for torrents.
func (c *client) TorrentVerify(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := torrentActionArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "torrent-verify", args)
	return err
}

// TorrentReannounce forces immediate tracker announce.
func (c *client) TorrentReannounce(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := torrentActionArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "torrent-reannounce", args)
	return err
}

// torrentGetArgs is used for torrent-get method.
type torrentGetArgs struct {
	IDs    any      `json:"ids,omitempty"`
	Fields []string `json:"fields"`
	Format string   `json:"format,omitempty"`
}

// TorrentGet retrieves information about torrents.
func (c *client) TorrentGet(ctx context.Context, fields []string, ids []int64) (*TorrentGetResult, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	args := torrentGetArgs{
		Fields: fields,
	}
	// Only set IDs if not empty (nil any vs any holding nil slice)
	if len(ids) > 0 {
		args.IDs = ids
	}

	resp, err := c.transport.Do(ctx, "torrent-get", args)
	if err != nil {
		return nil, err
	}

	var result TorrentGetResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return &result, nil
}

// TorrentGetByHash retrieves torrents by hash strings.
func (c *client) TorrentGetByHash(ctx context.Context, fields, hashes []string) (*TorrentGetResult, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	args := torrentGetArgs{
		IDs:    hashes,
		Fields: fields,
	}

	resp, err := c.transport.Do(ctx, "torrent-get", args)
	if err != nil {
		return nil, err
	}

	var result TorrentGetResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return &result, nil
}

// TorrentGetRecentlyActive retrieves recently active torrents.
func (c *client) TorrentGetRecentlyActive(ctx context.Context, fields []string) (*TorrentGetResult, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	args := torrentGetArgs{
		IDs:    "recently-active",
		Fields: fields,
	}

	resp, err := c.transport.Do(ctx, "torrent-get", args)
	if err != nil {
		return nil, err
	}

	var result TorrentGetResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return &result, nil
}

// torrentSetArgs wraps TorrentSetArgs with IDs.
type torrentSetArgs struct {
	*TorrentSetArgs

	IDs any `json:"ids"`
}

// TorrentSet modifies properties of one or more torrents.
func (c *client) TorrentSet(ctx context.Context, ids []int64, args *TorrentSetArgs) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	fullArgs := torrentSetArgs{
		IDs:            ids,
		TorrentSetArgs: args,
	}

	_, err := c.transport.Do(ctx, "torrent-set", fullArgs)
	return err
}

// TorrentAdd adds a new torrent.
func (c *client) TorrentAdd(ctx context.Context, args *TorrentAddArgs) (*TorrentAddResult, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	resp, err := c.transport.Do(ctx, "torrent-add", args)
	if err != nil {
		return nil, err
	}

	var result TorrentAddResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return &result, nil
}

// torrentRemoveArgs is used for torrent-remove method.
type torrentRemoveArgs struct {
	IDs             any  `json:"ids"`
	DeleteLocalData bool `json:"delete-local-data,omitempty"`
}

// TorrentRemove removes torrents.
func (c *client) TorrentRemove(ctx context.Context, ids []int64, deleteLocalData bool) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := torrentRemoveArgs{
		IDs:             ids,
		DeleteLocalData: deleteLocalData,
	}

	_, err := c.transport.Do(ctx, "torrent-remove", args)
	return err
}

// torrentSetLocationArgs is used for torrent-set-location method.
type torrentSetLocationArgs struct {
	IDs      any    `json:"ids"`
	Location string `json:"location"`
	Move     bool   `json:"move,omitempty"`
}

// TorrentSetLocation moves torrent data to a new location.
func (c *client) TorrentSetLocation(ctx context.Context, ids []int64, location string, move bool) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := torrentSetLocationArgs{
		IDs:      ids,
		Location: location,
		Move:     move,
	}

	_, err := c.transport.Do(ctx, "torrent-set-location", args)
	return err
}

// torrentRenamePathArgs is used for torrent-rename-path method.
type torrentRenamePathArgs struct {
	IDs  int64  `json:"ids"`
	Path string `json:"path"`
	Name string `json:"name"`
}

// TorrentRenamePath renames a file or directory within a torrent.
func (c *client) TorrentRenamePath(ctx context.Context, torrentID int64, path, name string) (*TorrentRenameResult, error) {
	if err := c.checkClosed(); err != nil {
		return nil, err
	}

	args := torrentRenamePathArgs{
		IDs:  torrentID,
		Path: path,
		Name: name,
	}

	resp, err := c.transport.Do(ctx, "torrent-rename-path", args)
	if err != nil {
		return nil, err
	}

	var result TorrentRenameResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, errors.Wrap(err, ErrUnmarshalResponse.Error())
	}

	return &result, nil
}
