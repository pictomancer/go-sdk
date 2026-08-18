package pictomancer

// OptimizeGeneratedParams tunes the optimize_generated operation. Zero values
// are omitted; the server defaults Format to webp and Strip to true, so Strip
// is a pointer (use Bool) to be able to send an explicit false.
type OptimizeGeneratedParams struct {
	Format        string         `json:"format,omitempty"`
	Q             int            `json:"q,omitempty"`
	QualityTarget float64        `json:"quality_target,omitempty"`
	MaxDimension  int            `json:"max_dimension,omitempty"`
	Strip         *bool          `json:"strip,omitempty"`
	Extra         map[string]any `json:"-"`
	Delivery      *Delivery      `json:"delivery,omitempty"`
}
