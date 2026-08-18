package pictomancer

// OpResult is the outcome of an image operation. Exactly one of Bytes and
// Receipt is populated: Bytes for inline delivery, Receipt (etag, sha256,
// bytes_written, ...) for put_url/callback deliveries. Quality accompanies
// either when the operation ran a quality_target search, nil otherwise.
type OpResult struct {
	Bytes   []byte
	Receipt map[string]any
	Quality *QualityReport
}
