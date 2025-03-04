package filesystem

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	tempDirectory  = "temp"
	maxSegmentSize = 10
	testDataDir    = "test_data"
)

func TestWrite(t *testing.T) {
	err := os.Mkdir(tempDirectory, 0777)
	require.NoError(t, err)
	defer func() {
		err := os.RemoveAll(tempDirectory)
		require.NoError(t, err)
	}()

	segment := NewSegment(tempDirectory, maxSegmentSize)

	now = func() time.Time {
		return time.Unix(1, 0)
	}

	err = segment.Write([]byte("12345"))
	require.NoError(t, err)
	err = segment.Write([]byte("67890"))
	require.NoError(t, err)

	stat, err := os.Stat(tempDirectory + "/wal_1000.wal")
	require.NoError(t, err)
	require.Equal(t, int64(10), stat.Size())

	now = func() time.Time {
		return time.Unix(2, 0)
	}
}

func TestLoadData(t *testing.T) {
	segment := NewSegment(testDataDir, maxSegmentSize)

	data, err := segment.LoadData()
	require.NoError(t, err)
	for _, d := range data {
		require.NotEmpty(t, d)
	}
}
