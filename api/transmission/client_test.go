package transmission

import (
	"context"
	_ "embed"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testMetainfo = "dGVzdA==" // base64 "test"

//go:embed testdata/responses/session-get.json
var sessionGetResponse string

//go:embed testdata/responses/session-get-fields.json
var sessionGetFieldsResponse string

//go:embed testdata/responses/session-stats.json
var sessionStatsResponse string

//go:embed testdata/responses/free-space.json
var freeSpaceResponse string

//go:embed testdata/responses/port-test.json
var portTestResponse string

//go:embed testdata/responses/group-get.json
var groupGetResponse string

//go:embed testdata/responses/torrent-add.json
var torrentAddResponse string

//go:embed testdata/responses/torrent-add-duplicate.json
var torrentAddDuplicateResponse string

//go:embed testdata/responses/torrent-get-all.json
var torrentGetAllResponse string

//go:embed testdata/responses/torrent-get-recently-active.json
var torrentGetRecentlyActiveResponse string

//go:embed testdata/responses/torrent-get-nonexistent.json
var torrentGetNonexistentResponse string

//go:embed testdata/responses/error-free-space-invalid-path.json
var errorFreeSpaceInvalidPathResponse string

//go:embed testdata/responses/error-blocklist-update.json
var errorBlocklistUpdateResponse string

//go:embed testdata/responses/error-torrent-add-invalid.json
var errorTorrentAddInvalidResponse string

func newTestClient(mock *RoundTripperMock) (Client, error) {
	httpClient := &http.Client{Transport: mock}
	return New("http://localhost:9091/transmission/rpc", WithHTTPClient(httpClient))
}

func TestSessionGet(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(sessionGetResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		session, err := client.SessionGet(context.Background(), nil)
		require.NoError(t, err)
		assert.NotNil(t, session.Version)
		assert.NotNil(t, session.RPCVersion)
	})

	t.Run("WithFields", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(sessionGetFieldsResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		session, err := client.SessionGet(context.Background(), []string{"version", "rpc-version"})
		require.NoError(t, err)
		assert.NotNil(t, session.Version)
		assert.NotNil(t, session.RPCVersion)
	})
}

func TestSessionStats(t *testing.T) {
	t.Parallel()

	mock := NewRoundTripperMock()
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return sessionIDResponse(), nil
	})
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return jsonResponse(sessionStatsResponse), nil
	})

	client, err := newTestClient(mock)
	require.NoError(t, err)

	stats, err := client.SessionStats(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, stats.TorrentCount, 0)
}

func TestFreeSpace(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(freeSpaceResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		space, err := client.FreeSpace(context.Background(), "/config")
		require.NoError(t, err)
		assert.NotEmpty(t, space.Path)
		assert.Positive(t, space.SizeBytes)
	})

	t.Run("InvalidPath", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(errorFreeSpaceInvalidPathResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		_, err = client.FreeSpace(context.Background(), "/nonexistent")
		assert.True(t, IsRPCError(err))
	})
}

func TestPortTest(t *testing.T) {
	t.Parallel()

	mock := NewRoundTripperMock()
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return sessionIDResponse(), nil
	})
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return jsonResponse(portTestResponse), nil
	})

	client, err := newTestClient(mock)
	require.NoError(t, err)

	result, err := client.PortTest(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestBlocklistUpdate(t *testing.T) {
	t.Parallel()

	mock := NewRoundTripperMock()
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return sessionIDResponse(), nil
	})
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return jsonResponse(errorBlocklistUpdateResponse), nil
	})

	client, err := newTestClient(mock)
	require.NoError(t, err)

	_, err = client.BlocklistUpdate(context.Background())
	assert.True(t, IsRPCError(err))
}

func TestGroupGet(t *testing.T) {
	t.Parallel()

	mock := NewRoundTripperMock()
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return sessionIDResponse(), nil
	})
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return jsonResponse(groupGetResponse), nil
	})

	client, err := newTestClient(mock)
	require.NoError(t, err)

	groups, err := client.GroupGet(context.Background(), nil)
	require.NoError(t, err)
	assert.NotNil(t, groups)
}

func TestTorrentAdd(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(torrentAddResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		metainfo := testMetainfo
		result, err := client.TorrentAdd(context.Background(), &TorrentAddArgs{
			Metainfo: &metainfo,
		})
		require.NoError(t, err)
		require.NotNil(t, result.TorrentAdded)
		assert.Equal(t, "ubuntu-24.04.3-live-server-amd64.iso", result.TorrentAdded.Name)
	})

	t.Run("Duplicate", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(torrentAddDuplicateResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		metainfo := testMetainfo
		result, err := client.TorrentAdd(context.Background(), &TorrentAddArgs{
			Metainfo: &metainfo,
		})
		require.NoError(t, err)
		assert.NotNil(t, result.TorrentDuplicate)
	})

	t.Run("Invalid", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(errorTorrentAddInvalidResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		metainfo := "bm90IGEgdG9ycmVudA==" // base64 "not a torrent"
		_, err = client.TorrentAdd(context.Background(), &TorrentAddArgs{
			Metainfo: &metainfo,
		})
		assert.True(t, IsRPCError(err))
	})
}

func TestTorrentGet(t *testing.T) {
	t.Parallel()

	t.Run("All", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(torrentGetAllResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		result, err := client.TorrentGet(context.Background(), []string{"id", "name"}, nil)
		require.NoError(t, err)
		assert.NotEmpty(t, result.Torrents)
	})

	t.Run("ByID", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			// Use torrent-get-all as it has actual torrent data
			return jsonResponse(torrentGetAllResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		result, err := client.TorrentGet(context.Background(), []string{"id", "name"}, []int64{1})
		require.NoError(t, err)
		assert.NotEmpty(t, result.Torrents)
	})

	t.Run("Nonexistent", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(torrentGetNonexistentResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		result, err := client.TorrentGet(context.Background(), []string{"id", "name"}, []int64{99999})
		require.NoError(t, err)
		assert.Empty(t, result.Torrents)
	})

	t.Run("RecentlyActive", func(t *testing.T) {
		t.Parallel()

		mock := NewRoundTripperMock()
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return sessionIDResponse(), nil
		})
		mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
			return jsonResponse(torrentGetRecentlyActiveResponse), nil
		})

		client, err := newTestClient(mock)
		require.NoError(t, err)

		result, err := client.TorrentGetRecentlyActive(context.Background(), []string{"id", "name"})
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestClientClose(t *testing.T) {
	t.Parallel()

	mock := NewRoundTripperMock()
	mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
		return sessionIDResponse(), nil
	})

	client, err := newTestClient(mock)
	require.NoError(t, err)

	err = client.Close()
	require.NoError(t, err)

	err = client.Close()
	assert.ErrorIs(t, err, ErrClosed)
}
