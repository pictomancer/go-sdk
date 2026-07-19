package pictomancer

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSourceHelpers(t *testing.T) {
	t.Parallel()

	imageBytes := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	wantSource := base64.StdEncoding.EncodeToString(imageBytes)

	t.Run("from bytes returns base64", func(t *testing.T) {
		t.Parallel()

		got := SourceFromBytes(imageBytes)

		require.Equal(t, wantSource, got)
	})

	t.Run("from reader reads and encodes", func(t *testing.T) {
		t.Parallel()

		got, err := SourceFromReader(bytes.NewReader(imageBytes))

		require.NoError(t, err)
		require.Equal(t, wantSource, got)
	})

	t.Run("from path reads and encodes", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join(t.TempDir(), "image.png")
		require.NoError(t, os.WriteFile(path, imageBytes, 0o600))

		got, err := SourceFromPath(path)

		require.NoError(t, err)
		require.Equal(t, wantSource, got)
	})

	t.Run("from path propagates missing file error", func(t *testing.T) {
		t.Parallel()

		_, err := SourceFromPath(filepath.Join(t.TempDir(), "missing.png"))

		require.Error(t, err)
	})
}
