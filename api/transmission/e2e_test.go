//go:build e2e

package transmission_test

import (
	"context"
	_ "embed"
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/lexfrei/go-transmission/api/transmission"
)

//go:embed testdata/ubuntu-24.04.3-live-server-amd64.iso.torrent
var ubuntuTorrent []byte

//go:embed testdata/Rocky-10.0-aarch64-dvd1.torrent
var rockyTorrent []byte

func setupTransmission(t *testing.T) (transmission.Client, func()) {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "lscr.io/linuxserver/transmission:latest",
		ExposedPorts: []string{"9091/tcp"},
		Env: map[string]string{
			"PUID": "1000",
			"PGID": "1000",
			"TZ":   "Etc/UTC",
		},
		WaitingFor: wait.ForHTTP("/transmission/web/").
			WithPort("9091/tcp").
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            false,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "9091")
	require.NoError(t, err)

	url := "http://" + host + ":" + port.Port() + "/transmission/rpc"

	client, err := transmission.New(url, transmission.WithTimeout(30*time.Second))
	require.NoError(t, err)

	cleanup := func() {
		_ = client.Close()
		_ = container.Terminate(ctx)
	}

	return client, cleanup
}

func TestE2E(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests in short mode")
	}

	client, cleanup := setupTransmission(t)
	defer cleanup()

	ctx := context.Background()

	// Store torrent IDs for later tests
	var torrentIDs []int64
	var torrentHashes []string

	t.Run("Session", func(t *testing.T) {
		t.Run("Get", func(t *testing.T) {
			session, err := client.SessionGet(ctx, nil)
			require.NoError(t, err)
			assert.NotNil(t, session.Version)
			assert.NotNil(t, session.RPCVersion)
		})

		t.Run("GetFields", func(t *testing.T) {
			session, err := client.SessionGet(ctx, []string{"version", "rpc-version"})
			require.NoError(t, err)
			assert.NotNil(t, session.Version)
			assert.NotNil(t, session.RPCVersion)
		})

		t.Run("Set", func(t *testing.T) {
			downloadLimit := int64(1234)
			err := client.SessionSet(ctx, &transmission.SessionSetArgs{
				SpeedLimitDown: &downloadLimit,
			})
			require.NoError(t, err)

			// Verify the setting was applied
			session, err := client.SessionGet(ctx, []string{"speed-limit-down"})
			require.NoError(t, err)
			require.NotNil(t, session.SpeedLimitDown)
			assert.Equal(t, int64(1234), *session.SpeedLimitDown)
		})

		t.Run("Stats", func(t *testing.T) {
			stats, err := client.SessionStats(ctx)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, stats.TorrentCount, 0)
		})
	})

	t.Run("System", func(t *testing.T) {
		t.Run("FreeSpace", func(t *testing.T) {
			space, err := client.FreeSpace(ctx, "/config")
			require.NoError(t, err)
			assert.NotEmpty(t, space.Path)
			assert.Greater(t, space.SizeBytes, int64(0))
		})

		t.Run("PortTest", func(t *testing.T) {
			// Port test may fail in container, just check no error
			_, err := client.PortTest(ctx)
			require.NoError(t, err)
		})

		t.Run("BlocklistUpdate", func(t *testing.T) {
			// Will fail with RPC error if no blocklist URL configured - that's expected
			_, err := client.BlocklistUpdate(ctx)
			// We just verify the method works, error is expected without blocklist URL
			assert.True(t, err == nil || transmission.IsRPCError(err))
		})
	})

	t.Run("BandwidthGroups", func(t *testing.T) {
		t.Run("Set", func(t *testing.T) {
			err := client.GroupSet(ctx, &transmission.BandwidthGroup{
				Name:                  "test-group",
				SpeedLimitUp:          int64Ptr(500),
				SpeedLimitDown:        int64Ptr(1000),
				SpeedLimitUpEnabled:   boolPtr(true),
				SpeedLimitDownEnabled: boolPtr(true),
				HonorsSessionLimits:   boolPtr(true),
			})
			require.NoError(t, err)
		})

		t.Run("Get", func(t *testing.T) {
			groups, err := client.GroupGet(ctx, nil)
			require.NoError(t, err)
			assert.NotEmpty(t, groups)

			// Find test-group and verify values
			var found bool
			for _, g := range groups {
				if g.Name == "test-group" {
					found = true
					require.NotNil(t, g.SpeedLimitUp)
					assert.Equal(t, int64(500), *g.SpeedLimitUp)
					require.NotNil(t, g.SpeedLimitDown)
					assert.Equal(t, int64(1000), *g.SpeedLimitDown)
					break
				}
			}
			assert.True(t, found, "test-group should exist")
		})

		t.Run("GetByName", func(t *testing.T) {
			groups, err := client.GroupGet(ctx, []string{"test-group"})
			require.NoError(t, err)
			require.Len(t, groups, 1)
			assert.Equal(t, "test-group", groups[0].Name)
		})
	})

	t.Run("TorrentAdd", func(t *testing.T) {
		t.Run("AddUbuntu", func(t *testing.T) {
			metainfo := base64.StdEncoding.EncodeToString(ubuntuTorrent)
			paused := true
			result, err := client.TorrentAdd(ctx, &transmission.TorrentAddArgs{
				Metainfo: &metainfo,
				Paused:   &paused,
			})
			require.NoError(t, err)
			require.NotNil(t, result.TorrentAdded)
			torrentIDs = append(torrentIDs, result.TorrentAdded.ID)
			torrentHashes = append(torrentHashes, result.TorrentAdded.HashString)
		})

		t.Run("AddRocky", func(t *testing.T) {
			metainfo := base64.StdEncoding.EncodeToString(rockyTorrent)
			paused := true
			result, err := client.TorrentAdd(ctx, &transmission.TorrentAddArgs{
				Metainfo: &metainfo,
				Paused:   &paused,
			})
			require.NoError(t, err)
			require.NotNil(t, result.TorrentAdded)
			torrentIDs = append(torrentIDs, result.TorrentAdded.ID)
			torrentHashes = append(torrentHashes, result.TorrentAdded.HashString)
		})

		t.Run("AddDuplicate", func(t *testing.T) {
			metainfo := base64.StdEncoding.EncodeToString(ubuntuTorrent)
			result, err := client.TorrentAdd(ctx, &transmission.TorrentAddArgs{
				Metainfo: &metainfo,
			})
			require.NoError(t, err)
			assert.NotNil(t, result.TorrentDuplicate)
		})
	})

	t.Run("TorrentGet", func(t *testing.T) {
		t.Run("GetAll", func(t *testing.T) {
			result, err := client.TorrentGet(ctx, []string{"id", "name", "status"}, nil)
			require.NoError(t, err)
			assert.Len(t, result.Torrents, 2)
		})

		t.Run("GetByID", func(t *testing.T) {
			result, err := client.TorrentGet(ctx, []string{"id", "name"}, torrentIDs[:1])
			require.NoError(t, err)
			assert.Len(t, result.Torrents, 1)
		})

		t.Run("GetByHash", func(t *testing.T) {
			result, err := client.TorrentGetByHash(ctx, []string{"id", "name"}, torrentHashes[:1])
			require.NoError(t, err)
			assert.Len(t, result.Torrents, 1)
		})

		t.Run("GetRecentlyActive", func(t *testing.T) {
			result, err := client.TorrentGetRecentlyActive(ctx, []string{"id", "name"})
			require.NoError(t, err)
			assert.NotNil(t, result)
		})
	})

	t.Run("TorrentSet", func(t *testing.T) {
		downloadLimit := int64(567)
		err := client.TorrentSet(ctx, torrentIDs[:1], &transmission.TorrentSetArgs{
			DownloadLimit:   &downloadLimit,
			DownloadLimited: boolPtr(true),
		})
		require.NoError(t, err)

		// Verify the setting was applied
		result, err := client.TorrentGet(ctx, []string{"id", "downloadLimit", "downloadLimited"}, torrentIDs[:1])
		require.NoError(t, err)
		require.Len(t, result.Torrents, 1)
		require.NotNil(t, result.Torrents[0].DownloadLimit)
		assert.Equal(t, int64(567), *result.Torrents[0].DownloadLimit)
		require.NotNil(t, result.Torrents[0].DownloadLimited)
		assert.True(t, *result.Torrents[0].DownloadLimited)
	})

	t.Run("TorrentActions", func(t *testing.T) {
		t.Run("Start", func(t *testing.T) {
			err := client.TorrentStart(ctx, torrentIDs)
			require.NoError(t, err)
		})

		t.Run("Stop", func(t *testing.T) {
			err := client.TorrentStop(ctx, torrentIDs)
			require.NoError(t, err)
		})

		t.Run("StartNow", func(t *testing.T) {
			err := client.TorrentStartNow(ctx, torrentIDs[:1])
			require.NoError(t, err)
		})

		t.Run("Verify", func(t *testing.T) {
			err := client.TorrentVerify(ctx, torrentIDs[:1])
			require.NoError(t, err)
		})

		t.Run("Reannounce", func(t *testing.T) {
			err := client.TorrentReannounce(ctx, torrentIDs[:1])
			require.NoError(t, err)
		})
	})

	t.Run("TorrentLocation", func(t *testing.T) {
		err := client.TorrentSetLocation(ctx, torrentIDs[:1], "/downloads/moved", false)
		require.NoError(t, err)

		// Verify location was changed
		result, err := client.TorrentGet(ctx, []string{"id", "downloadDir"}, torrentIDs[:1])
		require.NoError(t, err)
		require.Len(t, result.Torrents, 1)
		require.NotNil(t, result.Torrents[0].DownloadDir)
		assert.Equal(t, "/downloads/moved", *result.Torrents[0].DownloadDir)
	})

	t.Run("TorrentRenamePath", func(t *testing.T) {
		// Get file info first
		result, err := client.TorrentGet(ctx, []string{"id", "name", "files"}, torrentIDs[1:2])
		require.NoError(t, err)
		require.Len(t, result.Torrents, 1)

		torrent := result.Torrents[0]
		if torrent.Name != nil {
			renamed, err := client.TorrentRenamePath(ctx, torrentIDs[1], *torrent.Name, "renamed-torrent")
			require.NoError(t, err)
			assert.Equal(t, "renamed-torrent", renamed.Name)
		}
	})

	t.Run("Queue", func(t *testing.T) {
		// Get initial queue positions
		getPositions := func() (int, int) {
			result, err := client.TorrentGet(ctx, []string{"id", "queuePosition"}, torrentIDs)
			require.NoError(t, err)
			require.Len(t, result.Torrents, 2)
			var pos0, pos1 int
			for _, tor := range result.Torrents {
				if tor.ID != nil && tor.QueuePosition != nil {
					if *tor.ID == torrentIDs[0] {
						pos0 = *tor.QueuePosition
					} else {
						pos1 = *tor.QueuePosition
					}
				}
			}
			return pos0, pos1
		}

		t.Run("MoveBottom", func(t *testing.T) {
			err := client.QueueMoveBottom(ctx, torrentIDs[:1])
			require.NoError(t, err)

			pos0, pos1 := getPositions()
			assert.Greater(t, pos0, pos1, "torrent 0 should be after torrent 1")
		})

		t.Run("MoveTop", func(t *testing.T) {
			err := client.QueueMoveTop(ctx, torrentIDs[:1])
			require.NoError(t, err)

			pos0, pos1 := getPositions()
			assert.Less(t, pos0, pos1, "torrent 0 should be before torrent 1")
		})

		t.Run("MoveUp", func(t *testing.T) {
			// First move to bottom, then move up
			_ = client.QueueMoveBottom(ctx, torrentIDs[:1])
			posBefore, _ := getPositions()

			err := client.QueueMoveUp(ctx, torrentIDs[:1])
			require.NoError(t, err)

			posAfter, _ := getPositions()
			assert.Less(t, posAfter, posBefore, "position should decrease after MoveUp")
		})

		t.Run("MoveDown", func(t *testing.T) {
			// First move to top, then move down
			_ = client.QueueMoveTop(ctx, torrentIDs[:1])
			posBefore, _ := getPositions()

			err := client.QueueMoveDown(ctx, torrentIDs[:1])
			require.NoError(t, err)

			posAfter, _ := getPositions()
			assert.Greater(t, posAfter, posBefore, "position should increase after MoveDown")
		})
	})

	t.Run("TorrentRemove", func(t *testing.T) {
		t.Run("RemoveWithoutData", func(t *testing.T) {
			// Verify torrent exists before removal
			result, err := client.TorrentGet(ctx, []string{"id"}, torrentIDs[:1])
			require.NoError(t, err)
			require.Len(t, result.Torrents, 1)

			err = client.TorrentRemove(ctx, torrentIDs[:1], false)
			require.NoError(t, err)

			// Verify torrent is gone
			result, err = client.TorrentGet(ctx, []string{"id"}, torrentIDs[:1])
			require.NoError(t, err)
			assert.Empty(t, result.Torrents, "torrent should be removed")
		})

		t.Run("RemoveWithData", func(t *testing.T) {
			// Verify torrent exists before removal
			result, err := client.TorrentGet(ctx, []string{"id"}, torrentIDs[1:])
			require.NoError(t, err)
			require.Len(t, result.Torrents, 1)

			err = client.TorrentRemove(ctx, torrentIDs[1:], true)
			require.NoError(t, err)

			// Verify torrent is gone
			result, err = client.TorrentGet(ctx, []string{"id"}, torrentIDs[1:])
			require.NoError(t, err)
			assert.Empty(t, result.Torrents, "torrent should be removed")
		})

		t.Run("VerifyAllRemoved", func(t *testing.T) {
			// Verify no torrents remain
			result, err := client.TorrentGet(ctx, []string{"id"}, nil)
			require.NoError(t, err)
			assert.Empty(t, result.Torrents, "all torrents should be removed")
		})
	})

	t.Run("ErrorCases", func(t *testing.T) {
		nonExistentID := int64(99999)

		t.Run("GetNonExistent", func(t *testing.T) {
			result, err := client.TorrentGet(ctx, []string{"id", "name"}, []int64{nonExistentID})
			require.NoError(t, err)
			assert.Empty(t, result.Torrents, "should return empty list for non-existent ID")
		})

		t.Run("SetNonExistent", func(t *testing.T) {
			limit := int64(100)
			err := client.TorrentSet(ctx, []int64{nonExistentID}, &transmission.TorrentSetArgs{
				DownloadLimit: &limit,
			})
			// Transmission doesn't error on non-existent IDs for set operations
			require.NoError(t, err)
		})

		t.Run("StartNonExistent", func(t *testing.T) {
			err := client.TorrentStart(ctx, []int64{nonExistentID})
			require.NoError(t, err) // Transmission doesn't error on non-existent IDs
		})

		t.Run("RemoveNonExistent", func(t *testing.T) {
			err := client.TorrentRemove(ctx, []int64{nonExistentID}, false)
			require.NoError(t, err) // Transmission doesn't error on non-existent IDs
		})

		t.Run("FreeSpaceInvalidPath", func(t *testing.T) {
			_, err := client.FreeSpace(ctx, "/nonexistent/path/that/does/not/exist")
			assert.True(t, transmission.IsRPCError(err), "should return RPC error for invalid path")
		})
	})

	t.Run("ClientClose", func(t *testing.T) {
		err := client.Close()
		require.NoError(t, err)

		// Second close should return error
		err = client.Close()
		assert.ErrorIs(t, err, transmission.ErrClosed)
	})
}

func boolPtr(b bool) *bool {
	return &b
}

func int64Ptr(i int64) *int64 {
	return &i
}
