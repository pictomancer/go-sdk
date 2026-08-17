package pictomancer

// CompressParams tunes the compress operation. Zero values are omitted.
// QualityTarget (0 < v <= 1) asks for the smallest file with SSIM >= target
// instead of a fixed Q; mutually exclusive with Q, requires an explicit
// Format (jpeg/webp/avif). The outcome lands in OpResult.Quality.
// Denoise (1-3), Equalize and Sharpen are opt-in enhancement modifiers.
type CompressParams struct {
	Format        string         `json:"format,omitempty"`
	Q             int            `json:"q,omitempty"`
	QualityTarget float64        `json:"quality_target,omitempty"`
	Strip         bool           `json:"strip,omitempty"`
	Autorot       bool           `json:"autorot,omitempty"`
	Denoise       int            `json:"denoise,omitempty"`
	Equalize      bool           `json:"equalize,omitempty"`
	Sharpen       bool           `json:"sharpen,omitempty"`
	Extra         map[string]any `json:"-"`
	Delivery      *Delivery      `json:"delivery,omitempty"`
}
