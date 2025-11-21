package transmission

import (
	"context"
	"sync"
)

// Client is the interface for interacting with Transmission RPC API.
// All methods accept a context for cancellation and timeout control.
type Client interface {
	// Torrent Actions

	// TorrentStart starts one or more torrents.
	// If ids is nil, starts all torrents.
	TorrentStart(ctx context.Context, ids []int64) error

	// TorrentStartNow starts torrents immediately, bypassing the queue.
	TorrentStartNow(ctx context.Context, ids []int64) error

	// TorrentStop stops one or more torrents.
	// If ids is nil, stops all torrents.
	TorrentStop(ctx context.Context, ids []int64) error

	// TorrentVerify verifies local data integrity for torrents.
	TorrentVerify(ctx context.Context, ids []int64) error

	// TorrentReannounce forces immediate tracker announce.
	TorrentReannounce(ctx context.Context, ids []int64) error

	// Torrent Accessors

	// TorrentGet retrieves information about torrents.
	// fields specifies which fields to return (required).
	// ids specifies which torrents (nil for all).
	// Use "recently-active" as a special value by passing RecentlyActive.
	TorrentGet(ctx context.Context, fields []string, ids []int64) (*TorrentGetResult, error)

	// TorrentGetByHash retrieves torrents by hash strings.
	TorrentGetByHash(ctx context.Context, fields, hashes []string) (*TorrentGetResult, error)

	// TorrentGetRecentlyActive retrieves recently active torrents.
	// Also returns IDs of torrents removed since last call.
	TorrentGetRecentlyActive(ctx context.Context, fields []string) (*TorrentGetResult, error)

	// Torrent Mutators

	// TorrentSet modifies properties of one or more torrents.
	TorrentSet(ctx context.Context, ids []int64, args *TorrentSetArgs) error

	// TorrentAdd adds a new torrent.
	// Returns information about the added torrent, or duplicate if already exists.
	TorrentAdd(ctx context.Context, args *TorrentAddArgs) (*TorrentAddResult, error)

	// TorrentRemove removes torrents.
	// If deleteLocalData is true, also deletes downloaded files.
	TorrentRemove(ctx context.Context, ids []int64, deleteLocalData bool) error

	// Torrent Operations

	// TorrentSetLocation moves torrent data to a new location.
	// If move is true, files are moved; otherwise just the path is updated.
	TorrentSetLocation(ctx context.Context, ids []int64, location string, move bool) error

	// TorrentRenamePath renames a file or directory within a torrent.
	// Only one torrent can be processed at a time.
	TorrentRenamePath(ctx context.Context, id int64, path, name string) (*TorrentRenameResult, error)

	// Session

	// SessionGet retrieves session configuration.
	// If fields is nil, returns all fields.
	SessionGet(ctx context.Context, fields []string) (*Session, error)

	// SessionSet modifies session configuration.
	SessionSet(ctx context.Context, args *SessionSetArgs) error

	// SessionStats retrieves session statistics.
	SessionStats(ctx context.Context) (*SessionStats, error)

	// SessionClose gracefully shuts down the Transmission daemon.
	SessionClose(ctx context.Context) error

	// Queue Management

	// QueueMoveTop moves torrents to the top of the queue.
	QueueMoveTop(ctx context.Context, ids []int64) error

	// QueueMoveUp moves torrents up one position in the queue.
	QueueMoveUp(ctx context.Context, ids []int64) error

	// QueueMoveDown moves torrents down one position in the queue.
	QueueMoveDown(ctx context.Context, ids []int64) error

	// QueueMoveBottom moves torrents to the bottom of the queue.
	QueueMoveBottom(ctx context.Context, ids []int64) error

	// System

	// BlocklistUpdate updates the IP blocklist.
	// Returns the number of blocked IP ranges.
	BlocklistUpdate(ctx context.Context) (int, error)

	// PortTest tests if the peer port is accessible.
	// Returns true if the port is open.
	PortTest(ctx context.Context) (bool, error)

	// FreeSpace checks available disk space at the specified path.
	FreeSpace(ctx context.Context, path string) (*FreeSpace, error)

	// Bandwidth Groups

	// GroupSet creates or modifies a bandwidth group.
	GroupSet(ctx context.Context, args *BandwidthGroup) error

	// GroupGet retrieves bandwidth group information.
	// If names is nil, returns all groups.
	GroupGet(ctx context.Context, names []string) ([]BandwidthGroup, error)

	// Close releases resources associated with the client.
	Close() error
}

// TorrentGetResult contains the result of TorrentGet.
type TorrentGetResult struct {
	Torrents []Torrent `json:"torrents"`
	// Removed contains IDs of torrents removed since last "recently-active" call.
	Removed []int64 `json:"removed,omitempty"`
}

// TorrentSetArgs contains parameters for TorrentSet.
type TorrentSetArgs struct {
	BandwidthPriority           *Priority      `json:"bandwidthPriority,omitempty"`
	DownloadLimit               *int64         `json:"downloadLimit,omitempty"`
	DownloadLimited             *bool          `json:"downloadLimited,omitempty"`
	FilesUnwanted               []int          `json:"files-unwanted,omitempty"`
	FilesWanted                 []int          `json:"files-wanted,omitempty"`
	Group                       *string        `json:"group,omitempty"`
	HonorsSessionLimits         *bool          `json:"honorsSessionLimits,omitempty"`
	Labels                      []string       `json:"labels,omitempty"`
	Location                    *string        `json:"location,omitempty"`
	PeerLimit                   *int           `json:"peer-limit,omitempty"`
	PriorityHigh                []int          `json:"priority-high,omitempty"`
	PriorityLow                 []int          `json:"priority-low,omitempty"`
	PriorityNormal              []int          `json:"priority-normal,omitempty"`
	QueuePosition               *int           `json:"queuePosition,omitempty"`
	SeedIdleLimit               *int64         `json:"seedIdleLimit,omitempty"`
	SeedIdleMode                *SeedIdleMode  `json:"seedIdleMode,omitempty"`
	SeedRatioLimit              *float64       `json:"seedRatioLimit,omitempty"`
	SeedRatioMode               *SeedRatioMode `json:"seedRatioMode,omitempty"`
	SequentialDownload          *bool          `json:"sequential_download,omitempty"`
	SequentialDownloadFromPiece *int           `json:"sequential_download_from_piece,omitempty"`
	TrackerList                 *string        `json:"trackerList,omitempty"`
	UploadLimit                 *int64         `json:"uploadLimit,omitempty"`
	UploadLimited               *bool          `json:"uploadLimited,omitempty"`
}

// TorrentAddArgs contains parameters for TorrentAdd.
type TorrentAddArgs struct {
	// Filename is a torrent file path, magnet link, or URL (required if Metainfo is empty).
	Filename *string `json:"filename,omitempty"`
	// Metainfo is base64-encoded .torrent content (required if Filename is empty).
	Metainfo *string `json:"metainfo,omitempty"`
	// DownloadDir is the download directory path.
	DownloadDir *string `json:"download-dir,omitempty"`
	// Cookies for URL downloads (format: "NAME=VALUE").
	Cookies *string `json:"cookies,omitempty"`
	// Paused starts the torrent paused if true.
	Paused *bool `json:"paused,omitempty"`
	// PeerLimit is the maximum number of peers.
	PeerLimit *int `json:"peer-limit,omitempty"`
	// BandwidthPriority sets the torrent priority.
	BandwidthPriority *Priority `json:"bandwidthPriority,omitempty"`
	// Labels are tags for the torrent.
	Labels []string `json:"labels,omitempty"`
	// FilesWanted are file indices to download.
	FilesWanted []int `json:"files-wanted,omitempty"`
	// FilesUnwanted are file indices to skip.
	FilesUnwanted []int `json:"files-unwanted,omitempty"`
	// PriorityHigh are file indices for high priority.
	PriorityHigh []int `json:"priority-high,omitempty"`
	// PriorityLow are file indices for low priority.
	PriorityLow []int `json:"priority-low,omitempty"`
	// PriorityNormal are file indices for normal priority.
	PriorityNormal []int `json:"priority-normal,omitempty"`
	// SequentialDownload enables sequential downloading.
	SequentialDownload *bool `json:"sequential_download,omitempty"`
	// SequentialDownloadFromPiece is the starting piece for sequential download.
	SequentialDownloadFromPiece *int `json:"sequential_download_from_piece,omitempty"`
}

// TorrentAddResult contains the result of TorrentAdd.
type TorrentAddResult struct {
	// TorrentAdded is set if a new torrent was added.
	TorrentAdded *TorrentAddedInfo `json:"torrent-added,omitempty"`
	// TorrentDuplicate is set if the torrent already exists.
	TorrentDuplicate *TorrentAddedInfo `json:"torrent-duplicate,omitempty"`
}

// TorrentRenameResult contains the result of TorrentRenamePath.
type TorrentRenameResult struct {
	Path string `json:"path"`
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

// SessionSetArgs contains parameters for SessionSet.
type SessionSetArgs struct {
	AltSpeedDown                     *int64          `json:"alt-speed-down,omitempty"`
	AltSpeedEnabled                  *bool           `json:"alt-speed-enabled,omitempty"`
	AltSpeedTimeBegin                *int            `json:"alt-speed-time-begin,omitempty"`
	AltSpeedTimeDay                  *int            `json:"alt-speed-time-day,omitempty"`
	AltSpeedTimeEnabled              *bool           `json:"alt-speed-time-enabled,omitempty"`
	AltSpeedTimeEnd                  *int            `json:"alt-speed-time-end,omitempty"`
	AltSpeedUp                       *int64          `json:"alt-speed-up,omitempty"`
	BlocklistEnabled                 *bool           `json:"blocklist-enabled,omitempty"`
	BlocklistURL                     *string         `json:"blocklist-url,omitempty"`
	CacheSizeMB                      *int            `json:"cache-size-mb,omitempty"`
	DefaultTrackers                  *string         `json:"default-trackers,omitempty"`
	DHTEnabled                       *bool           `json:"dht-enabled,omitempty"`
	DownloadDir                      *string         `json:"download-dir,omitempty"`
	DownloadQueueEnabled             *bool           `json:"download-queue-enabled,omitempty"`
	DownloadQueueSize                *int            `json:"download-queue-size,omitempty"`
	Encryption                       *EncryptionMode `json:"encryption,omitempty"`
	IdleSeedingLimit                 *int            `json:"idle-seeding-limit,omitempty"`
	IdleSeedingLimitEnabled          *bool           `json:"idle-seeding-limit-enabled,omitempty"`
	IncompleteDir                    *string         `json:"incomplete-dir,omitempty"`
	IncompleteDirEnabled             *bool           `json:"incomplete-dir-enabled,omitempty"`
	LPDEnabled                       *bool           `json:"lpd-enabled,omitempty"`
	PeerLimitGlobal                  *int            `json:"peer-limit-global,omitempty"`
	PeerLimitPerTorrent              *int            `json:"peer-limit-per-torrent,omitempty"`
	PeerPort                         *int            `json:"peer-port,omitempty"`
	PeerPortRandomOnStart            *bool           `json:"peer-port-random-on-start,omitempty"`
	PEXEnabled                       *bool           `json:"pex-enabled,omitempty"`
	PortForwardingEnabled            *bool           `json:"port-forwarding-enabled,omitempty"`
	PreferredTransports              []string        `json:"preferred_transports,omitempty"`
	QueueStalledEnabled              *bool           `json:"queue-stalled-enabled,omitempty"`
	QueueStalledMinutes              *int            `json:"queue-stalled-minutes,omitempty"`
	RenamePartialFiles               *bool           `json:"rename-partial-files,omitempty"`
	ScriptTorrentAddedEnabled        *bool           `json:"script-torrent-added-enabled,omitempty"`
	ScriptTorrentAddedFilename       *string         `json:"script-torrent-added-filename,omitempty"`
	ScriptTorrentDoneEnabled         *bool           `json:"script-torrent-done-enabled,omitempty"`
	ScriptTorrentDoneFilename        *string         `json:"script-torrent-done-filename,omitempty"`
	ScriptTorrentDoneSeedingEnabled  *bool           `json:"script-torrent-done-seeding-enabled,omitempty"`
	ScriptTorrentDoneSeedingFilename *string         `json:"script-torrent-done-seeding-filename,omitempty"`
	SeedQueueEnabled                 *bool           `json:"seed-queue-enabled,omitempty"`
	SeedQueueSize                    *int            `json:"seed-queue-size,omitempty"`
	SeedRatioLimit                   *float64        `json:"seedRatioLimit,omitempty"`
	SeedRatioLimited                 *bool           `json:"seedRatioLimited,omitempty"`
	SequentialDownload               *bool           `json:"sequential_download,omitempty"`
	SpeedLimitDown                   *int64          `json:"speed-limit-down,omitempty"`
	SpeedLimitDownEnabled            *bool           `json:"speed-limit-down-enabled,omitempty"`
	SpeedLimitUp                     *int64          `json:"speed-limit-up,omitempty"`
	SpeedLimitUpEnabled              *bool           `json:"speed-limit-up-enabled,omitempty"`
	StartAddedTorrents               *bool           `json:"start-added-torrents,omitempty"`
	TrashOriginalTorrentFiles        *bool           `json:"trash-original-torrent-files,omitempty"`
}

// client is the default implementation of Client.
type client struct {
	transport *httpTransport
	closed    bool
	mu        sync.RWMutex
}

// New creates a new Transmission RPC client.
func New(url string, opts ...Option) (Client, error) {
	cfg := &config{
		url: url,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.url == "" {
		return nil, ErrInvalidURL
	}

	transport, err := newHTTPTransport(cfg)
	if err != nil {
		return nil, err
	}

	return &client{
		transport: transport,
	}, nil
}

// Close releases resources associated with the client.
func (c *client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return ErrClosed
	}

	c.closed = true
	return nil
}

// checkClosed returns an error if the client is closed.
func (c *client) checkClosed() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return ErrClosed
	}
	return nil
}
