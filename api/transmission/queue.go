package transmission

import "context"

// queueMoveArgs is used for queue movement methods.
type queueMoveArgs struct {
	IDs any `json:"ids"`
}

// QueueMoveTop moves torrents to the top of the queue.
func (c *client) QueueMoveTop(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := queueMoveArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "queue-move-top", args)
	return err
}

// QueueMoveUp moves torrents up one position in the queue.
func (c *client) QueueMoveUp(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := queueMoveArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "queue-move-up", args)
	return err
}

// QueueMoveDown moves torrents down one position in the queue.
func (c *client) QueueMoveDown(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := queueMoveArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "queue-move-down", args)
	return err
}

// QueueMoveBottom moves torrents to the bottom of the queue.
func (c *client) QueueMoveBottom(ctx context.Context, ids []int64) error {
	if err := c.checkClosed(); err != nil {
		return err
	}

	args := queueMoveArgs{IDs: ids}
	_, err := c.transport.Do(ctx, "queue-move-bottom", args)
	return err
}
