package pictomancer

// OpResult is the outcome of an image operation. Exactly one field is
// populated: Bytes for inline delivery, Receipt (etag, sha256,
// bytes_written, ...) for put_url/callback deliveries.
type OpResult struct {
	Bytes   []byte
	Receipt map[string]any
}
