package pictomancer

// CompressParams tunes the compress operation. Zero values are omitted.
// QualityTarget (0 < v <= 1) asks for the smallest file with SSIM >= target
// instead of a fixed Q; mutually exclusive with Q, requires an explicit
// Format (jpeg/webp/avif). The outcome lands in OpResult.Quality.
type CompressParams struct {
	Format        string         `json:"format,omitempty"`
	Q             int            `json:"q,omitempty"`
	QualityTarget float64        `json:"quality_target,omitempty"`
	Strip         bool           `json:"strip,omitempty"`
	Extra         map[string]any `json:"-"`
	Delivery      *Delivery      `json:"delivery,omitempty"`
}
