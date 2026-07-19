package pictomancer

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// SourceFromBytes encodes in-memory image bytes as the raw base64
// source string the API accepts.
func SourceFromBytes(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// SourceFromReader reads the whole image from r and encodes it.
func SourceFromReader(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("read source: %w", err)
	}
	return SourceFromBytes(data), nil
}

// SourceFromPath reads a local image file and encodes it.
func SourceFromPath(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read source file: %w", err)
	}
	return SourceFromBytes(data), nil
}
