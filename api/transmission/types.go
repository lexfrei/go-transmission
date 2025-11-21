package transmission

// TorrentStatus represents the current state of a torrent.
type TorrentStatus int

// Torrent status values.
const (
	TorrentStatusStopped      TorrentStatus = 0 // Torrent is stopped
	TorrentStatusCheckWait    TorrentStatus = 1 // Queued to check files
	TorrentStatusCheck        TorrentStatus = 2 // Checking files
	TorrentStatusDownloadWait TorrentStatus = 3 // Queued to download
	TorrentStatusDownload     TorrentStatus = 4 // Downloading
	TorrentStatusSeedWait     TorrentStatus = 5 // Queued to seed
	TorrentStatusSeed         TorrentStatus = 6 // Seeding
)

// String returns a human-readable representation of the torrent status.
func (s TorrentStatus) String() string {
	switch s {
	case TorrentStatusStopped:
		return "Stopped"
	case TorrentStatusCheckWait:
		return "Check Wait"
	case TorrentStatusCheck:
		return "Checking"
	case TorrentStatusDownloadWait:
		return "Download Wait"
	case TorrentStatusDownload:
		return "Downloading"
	case TorrentStatusSeedWait:
		return "Seed Wait"
	case TorrentStatusSeed:
		return "Seeding"
	default:
		return "Unknown"
	}
}

// Priority represents file or torrent priority.
type Priority int

// Priority values.
const (
	PriorityLow    Priority = -1
	PriorityNormal Priority = 0
	PriorityHigh   Priority = 1
)

// EncryptionMode represents the peer encryption mode.
type EncryptionMode string

// Encryption mode values.
const (
	EncryptionRequired  EncryptionMode = "required"
	EncryptionPreferred EncryptionMode = "preferred"
	EncryptionTolerated EncryptionMode = "tolerated"
)

// SeedRatioMode represents the seed ratio limit mode.
type SeedRatioMode int

// Seed ratio mode values.
const (
	SeedRatioModeGlobal    SeedRatioMode = 0 // Use global settings
	SeedRatioModeSingle    SeedRatioMode = 1 // Use torrent-specific limit
	SeedRatioModeUnlimited SeedRatioMode = 2 // No limit
)

// SeedIdleMode represents the seed idle limit mode.
type SeedIdleMode int

// Seed idle mode values.
const (
	SeedIdleModeGlobal    SeedIdleMode = 0 // Use global settings
	SeedIdleModeSingle    SeedIdleMode = 1 // Use torrent-specific limit
	SeedIdleModeUnlimited SeedIdleMode = 2 // No limit
)

// Torrent represents a torrent in Transmission.
// Fields are optional and only populated if requested in TorrentGet.
type Torrent struct {
	// Identification
	ID         *int64  `json:"id,omitempty"`
	HashString *string `json:"hashString,omitempty"`
	Name       *string `json:"name,omitempty"`

	// Status
	Status      *TorrentStatus `json:"status,omitempty"`
	Error       *int           `json:"error,omitempty"`
	ErrorString *string        `json:"errorString,omitempty"`
	IsFinished  *bool          `json:"isFinished,omitempty"`
	IsStalled   *bool          `json:"isStalled,omitempty"`

	// Timestamps (Unix time)
	ActivityDate *int64 `json:"activityDate,omitempty"`
	AddedDate    *int64 `json:"addedDate,omitempty"`
	DoneDate     *int64 `json:"doneDate,omitempty"`
	EditDate     *int64 `json:"editDate,omitempty"`
	StartDate    *int64 `json:"startDate,omitempty"`
	DateCreated  *int64 `json:"dateCreated,omitempty"`

	// Progress (0.0 - 1.0)
	PercentComplete         *float64 `json:"percentComplete,omitempty"`
	PercentDone             *float64 `json:"percentDone,omitempty"`
	MetadataPercentComplete *float64 `json:"metadataPercentComplete,omitempty"`
	RecheckProgress         *float64 `json:"recheckProgress,omitempty"`

	// Speed (bytes per second)
	RateDownload *int64 `json:"rateDownload,omitempty"`
	RateUpload   *int64 `json:"rateUpload,omitempty"`

	// Time estimates (seconds)
	ETA     *int64 `json:"eta,omitempty"`
	ETAIdle *int64 `json:"etaIdle,omitempty"`

	// Duration (seconds)
	SecondsDownloading *int64 `json:"secondsDownloading,omitempty"`
	SecondsSeeding     *int64 `json:"secondsSeeding,omitempty"`

	// Size (bytes)
	TotalSize        *int64 `json:"totalSize,omitempty"`
	SizeWhenDone     *int64 `json:"sizeWhenDone,omitempty"`
	LeftUntilDone    *int64 `json:"leftUntilDone,omitempty"`
	DesiredAvailable *int64 `json:"desiredAvailable,omitempty"`

	// Transfer totals (bytes)
	DownloadedEver *int64 `json:"downloadedEver,omitempty"`
	UploadedEver   *int64 `json:"uploadedEver,omitempty"`
	CorruptEver    *int64 `json:"corruptEver,omitempty"`
	HaveValid      *int64 `json:"haveValid,omitempty"`
	HaveUnchecked  *int64 `json:"haveUnchecked,omitempty"`

	// Bandwidth settings
	BandwidthPriority   *Priority `json:"bandwidthPriority,omitempty"`
	DownloadLimit       *int64    `json:"downloadLimit,omitempty"`
	DownloadLimited     *bool     `json:"downloadLimited,omitempty"`
	UploadLimit         *int64    `json:"uploadLimit,omitempty"`
	UploadLimited       *bool     `json:"uploadLimited,omitempty"`
	HonorsSessionLimits *bool     `json:"honorsSessionLimits,omitempty"`
	Group               *string   `json:"group,omitempty"`

	// Seed limits
	SeedRatioLimit *float64       `json:"seedRatioLimit,omitempty"`
	SeedRatioMode  *SeedRatioMode `json:"seedRatioMode,omitempty"`
	SeedIdleLimit  *int64         `json:"seedIdleLimit,omitempty"`
	SeedIdleMode   *SeedIdleMode  `json:"seedIdleMode,omitempty"`

	// Peers
	PeersConnected     *int       `json:"peersConnected,omitempty"`
	PeersGettingFromUs *int       `json:"peersGettingFromUs,omitempty"`
	PeersSendingToUs   *int       `json:"peersSendingToUs,omitempty"`
	PeerLimit          *int       `json:"peer-limit,omitempty"`
	MaxConnectedPeers  *int       `json:"maxConnectedPeers,omitempty"`
	Peers              []Peer     `json:"peers,omitempty"`
	PeersFrom          *PeersFrom `json:"peersFrom,omitempty"`

	// Files
	FileCount  *int       `json:"fileCount,omitempty"`
	Files      []File     `json:"files,omitempty"`
	FileStats  []FileStat `json:"fileStats,omitempty"`
	Priorities []Priority `json:"priorities,omitempty"`
	Wanted     []bool     `json:"wanted,omitempty"`

	// Trackers
	Trackers     []Tracker     `json:"trackers,omitempty"`
	TrackerStats []TrackerStat `json:"trackerStats,omitempty"`
	TrackerList  *string       `json:"trackerList,omitempty"`

	// Pieces
	Pieces       *string   `json:"pieces,omitempty"` // Base64-encoded bitfield
	PieceCount   *int      `json:"pieceCount,omitempty"`
	PieceSize    *int64    `json:"pieceSize,omitempty"`
	Availability []float64 `json:"availability,omitempty"`

	// Metadata
	Comment     *string  `json:"comment,omitempty"`
	Creator     *string  `json:"creator,omitempty"`
	IsPrivate   *bool    `json:"isPrivate,omitempty"`
	Labels      []string `json:"labels,omitempty"`
	MagnetLink  *string  `json:"magnetLink,omitempty"`
	TorrentFile *string  `json:"torrentFile,omitempty"`

	// Location
	DownloadDir *string `json:"downloadDir,omitempty"`

	// Queue
	QueuePosition *int `json:"queuePosition,omitempty"`

	// Sequential download
	SequentialDownload          *bool `json:"sequential_download,omitempty"`
	SequentialDownloadFromPiece *int  `json:"sequential_download_from_piece,omitempty"`

	// Other
	PrimaryMimeType *string  `json:"primary-mime-type,omitempty"`
	Webseeds        []string `json:"webseeds,omitempty"`
}

// File represents a file within a torrent.
type File struct {
	Name           string `json:"name"`
	Length         int64  `json:"length"`
	BytesCompleted int64  `json:"bytesCompleted"`
	BeginPiece     *int   `json:"begin_piece,omitempty"`
	EndPiece       *int   `json:"end_piece,omitempty"`
}

// FileStat represents the status of a file within a torrent.
type FileStat struct {
	BytesCompleted int64    `json:"bytesCompleted"`
	Wanted         bool     `json:"wanted"`
	Priority       Priority `json:"priority"`
}

// Peer represents a connected peer.
type Peer struct {
	Address            string  `json:"address"`
	ClientName         string  `json:"clientName"`
	ClientIsChoked     bool    `json:"clientIsChoked"`
	ClientIsInterested bool    `json:"clientIsInterested"`
	FlagStr            string  `json:"flagStr"`
	IsDownloadingFrom  bool    `json:"isDownloadingFrom"`
	IsEncrypted        bool    `json:"isEncrypted"`
	IsIncoming         bool    `json:"isIncoming"`
	IsUploadingTo      bool    `json:"isUploadingTo"`
	IsUTP              bool    `json:"isUTP"`
	PeerIsChoked       bool    `json:"peerIsChoked"`
	PeerIsInterested   bool    `json:"peerIsInterested"`
	Port               int     `json:"port"`
	Progress           float64 `json:"progress"`
	RateToClient       int64   `json:"rateToClient"`
	RateToPeer         int64   `json:"rateToPeer"`
}

// PeersFrom indicates where peers were discovered from.
type PeersFrom struct {
	FromCache    int `json:"fromCache"`
	FromDHT      int `json:"fromDht"`
	FromIncoming int `json:"fromIncoming"`
	FromLPD      int `json:"fromLpd"`
	FromLTEP     int `json:"fromLtep"`
	FromPEX      int `json:"fromPex"`
	FromTracker  int `json:"fromTracker"`
}

// Tracker represents a tracker for a torrent.
type Tracker struct {
	Announce string `json:"announce"`
	ID       int    `json:"id"`
	Scrape   string `json:"scrape,omitempty"`
	Sitename string `json:"sitename,omitempty"`
	Tier     int    `json:"tier"`
}

// TrackerStat represents tracker statistics.
type TrackerStat struct {
	Announce              string `json:"announce"`
	AnnounceState         int    `json:"announceState"`
	DownloadCount         int    `json:"downloadCount"`
	HasAnnounced          bool   `json:"hasAnnounced"`
	HasScraped            bool   `json:"hasScraped"`
	Host                  string `json:"host"`
	ID                    int    `json:"id"`
	IsBackup              bool   `json:"isBackup"`
	LastAnnouncePeerCount int    `json:"lastAnnouncePeerCount"`
	LastAnnounceResult    string `json:"lastAnnounceResult"`
	LastAnnounceStartTime int64  `json:"lastAnnounceStartTime"`
	LastAnnounceSucceeded bool   `json:"lastAnnounceSucceeded"`
	LastAnnounceTime      int64  `json:"lastAnnounceTime"`
	LastAnnounceTimedOut  bool   `json:"lastAnnounceTimedOut"`
	LastScrapeResult      string `json:"lastScrapeResult"`
	LastScrapeStartTime   int64  `json:"lastScrapeStartTime"`
	LastScrapeSucceeded   bool   `json:"lastScrapeSucceeded"`
	LastScrapeTime        int64  `json:"lastScrapeTime"`
	LastScrapeTimedOut    bool   `json:"lastScrapeTimedOut"`
	LeecherCount          int    `json:"leecherCount"`
	NextAnnounceTime      int64  `json:"nextAnnounceTime"`
	NextScrapeTime        int64  `json:"nextScrapeTime"`
	Scrape                string `json:"scrape"`
	ScrapeState           int    `json:"scrapeState"`
	SeederCount           int    `json:"seederCount"`
	Sitename              string `json:"sitename"`
	Tier                  int    `json:"tier"`
}

// TorrentAddedInfo contains information about a newly added torrent.
type TorrentAddedInfo struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	HashString string `json:"hashString"`
}

// Session represents Transmission session configuration.
type Session struct {
	// Speed limits
	AltSpeedDown          *int64 `json:"alt-speed-down,omitempty"`
	AltSpeedEnabled       *bool  `json:"alt-speed-enabled,omitempty"`
	AltSpeedTimeBegin     *int   `json:"alt-speed-time-begin,omitempty"`
	AltSpeedTimeDay       *int   `json:"alt-speed-time-day,omitempty"`
	AltSpeedTimeEnabled   *bool  `json:"alt-speed-time-enabled,omitempty"`
	AltSpeedTimeEnd       *int   `json:"alt-speed-time-end,omitempty"`
	AltSpeedUp            *int64 `json:"alt-speed-up,omitempty"`
	SpeedLimitDown        *int64 `json:"speed-limit-down,omitempty"`
	SpeedLimitDownEnabled *bool  `json:"speed-limit-down-enabled,omitempty"`
	SpeedLimitUp          *int64 `json:"speed-limit-up,omitempty"`
	SpeedLimitUpEnabled   *bool  `json:"speed-limit-up-enabled,omitempty"`

	// Blocklist
	BlocklistEnabled *bool   `json:"blocklist-enabled,omitempty"`
	BlocklistSize    *int    `json:"blocklist-size,omitempty"` // Read-only
	BlocklistURL     *string `json:"blocklist-url,omitempty"`

	// Cache
	CacheSizeMB *int `json:"cache-size-mb,omitempty"`

	// DHT, PEX, LPD
	DHTEnabled *bool `json:"dht-enabled,omitempty"`
	PEXEnabled *bool `json:"pex-enabled,omitempty"`
	LPDEnabled *bool `json:"lpd-enabled,omitempty"`

	// Directories
	DownloadDir          *string `json:"download-dir,omitempty"`
	IncompleteDir        *string `json:"incomplete-dir,omitempty"`
	IncompleteDirEnabled *bool   `json:"incomplete-dir-enabled,omitempty"`
	ConfigDir            *string `json:"config-dir,omitempty"` // Read-only

	// Encryption
	Encryption *EncryptionMode `json:"encryption,omitempty"`

	// Peer settings
	PeerLimitGlobal       *int  `json:"peer-limit-global,omitempty"`
	PeerLimitPerTorrent   *int  `json:"peer-limit-per-torrent,omitempty"`
	PeerPort              *int  `json:"peer-port,omitempty"`
	PeerPortRandomOnStart *bool `json:"peer-port-random-on-start,omitempty"`
	PortForwardingEnabled *bool `json:"port-forwarding-enabled,omitempty"`

	// Queue
	DownloadQueueEnabled *bool `json:"download-queue-enabled,omitempty"`
	DownloadQueueSize    *int  `json:"download-queue-size,omitempty"`
	SeedQueueEnabled     *bool `json:"seed-queue-enabled,omitempty"`
	SeedQueueSize        *int  `json:"seed-queue-size,omitempty"`
	QueueStalledEnabled  *bool `json:"queue-stalled-enabled,omitempty"`
	QueueStalledMinutes  *int  `json:"queue-stalled-minutes,omitempty"`

	// Seed limits
	IdleSeedingLimit        *int     `json:"idle-seeding-limit,omitempty"`
	IdleSeedingLimitEnabled *bool    `json:"idle-seeding-limit-enabled,omitempty"`
	SeedRatioLimit          *float64 `json:"seedRatioLimit,omitempty"`
	SeedRatioLimited        *bool    `json:"seedRatioLimited,omitempty"`

	// Scripts
	ScriptTorrentAddedEnabled        *bool   `json:"script-torrent-added-enabled,omitempty"`
	ScriptTorrentAddedFilename       *string `json:"script-torrent-added-filename,omitempty"`
	ScriptTorrentDoneEnabled         *bool   `json:"script-torrent-done-enabled,omitempty"`
	ScriptTorrentDoneFilename        *string `json:"script-torrent-done-filename,omitempty"`
	ScriptTorrentDoneSeedingEnabled  *bool   `json:"script-torrent-done-seeding-enabled,omitempty"`
	ScriptTorrentDoneSeedingFilename *string `json:"script-torrent-done-seeding-filename,omitempty"`

	// Misc settings
	DefaultTrackers           *string  `json:"default-trackers,omitempty"`
	RenamePartialFiles        *bool    `json:"rename-partial-files,omitempty"`
	StartAddedTorrents        *bool    `json:"start-added-torrents,omitempty"`
	TrashOriginalTorrentFiles *bool    `json:"trash-original-torrent-files,omitempty"`
	SequentialDownload        *bool    `json:"sequential_download,omitempty"`
	PreferredTransports       []string `json:"preferred_transports,omitempty"`

	// Version info (read-only)
	RPCVersion        *int    `json:"rpc-version,omitempty"`
	RPCVersionMinimum *int    `json:"rpc-version-minimum,omitempty"`
	RPCVersionSemver  *string `json:"rpc-version-semver,omitempty"`
	Version           *string `json:"version,omitempty"`
	SessionID         *string `json:"session-id,omitempty"`

	// Units (read-only)
	Units *Units `json:"units,omitempty"`
}

// Units represents the measurement units used by Transmission.
type Units struct {
	SpeedUnits  []string `json:"speed-units"`
	SpeedBytes  int      `json:"speed-bytes"`
	SizeUnits   []string `json:"size-units"`
	SizeBytes   int      `json:"size-bytes"`
	MemoryUnits []string `json:"memory-units"`
	MemoryBytes int      `json:"memory-bytes"`
}

// SessionStats represents session statistics.
type SessionStats struct {
	ActiveTorrentCount int   `json:"activeTorrentCount"`
	DownloadSpeed      int64 `json:"downloadSpeed"`
	PausedTorrentCount int   `json:"pausedTorrentCount"`
	TorrentCount       int   `json:"torrentCount"`
	UploadSpeed        int64 `json:"uploadSpeed"`
	CumulativeStats    Stats `json:"cumulative-stats"`
	CurrentStats       Stats `json:"current-stats"`
}

// Stats represents transfer statistics.
type Stats struct {
	UploadedBytes   int64 `json:"uploadedBytes"`
	DownloadedBytes int64 `json:"downloadedBytes"`
	FilesAdded      int   `json:"filesAdded"`
	SessionCount    int   `json:"sessionCount"`
	SecondsActive   int64 `json:"secondsActive"`
}

// FreeSpace represents free space information for a path.
type FreeSpace struct {
	Path      string `json:"path"`
	SizeBytes int64  `json:"size-bytes"`
	TotalSize *int64 `json:"total_size,omitempty"`
}

// BandwidthGroup represents a bandwidth group configuration.
type BandwidthGroup struct {
	Name                  string `json:"name"`
	HonorsSessionLimits   *bool  `json:"honorsSessionLimits,omitempty"`
	SpeedLimitDownEnabled *bool  `json:"speed-limit-down-enabled,omitempty"`
	SpeedLimitDown        *int64 `json:"speed-limit-down,omitempty"`
	SpeedLimitUpEnabled   *bool  `json:"speed-limit-up-enabled,omitempty"`
	SpeedLimitUp          *int64 `json:"speed-limit-up,omitempty"`
}
