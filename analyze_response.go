package pictomancer

// AnalyzeResponse is the metadata of a fetched image. Always free.
type AnalyzeResponse struct {
	SizeBytes int64 `json:"size_bytes"`
}
