package pictomancer

// ConvertParams tunes the convert operation. Zero values are omitted.
// Effort is a pointer because 0 is a valid AVIF encoder value distinct
// from unset (API default 2); use Int to build it.
// QualityTarget (0 < v <= 1) asks for the smallest file with SSIM >= target
// instead of a fixed Q; mutually exclusive with Q and Lossless, only for
// jpeg/webp/avif targets. The outcome lands in OpResult.Quality.
type ConvertParams struct {
	Q             int            `json:"q,omitempty"`
	QualityTarget float64        `json:"quality_target,omitempty"`
	Strip         bool           `json:"strip,omitempty"`
	Lossless      bool           `json:"lossless,omitempty"`
	Effort        *int           `json:"effort,omitempty"`
	Autorot       bool           `json:"autorot,omitempty"`
	Extra         map[string]any `json:"-"`
	Delivery      *Delivery      `json:"delivery,omitempty"`
}
