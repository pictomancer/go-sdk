package pictomancer

// CropParams tunes the crop operation. Zero values are omitted.
type CropParams struct {
	Format   string         `json:"format,omitempty"`
	Extra    map[string]any `json:"-"`
	Delivery *Delivery      `json:"delivery,omitempty"`
}
