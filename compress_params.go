package pictomancer

// CompressParams tunes the compress operation. Zero values are omitted.
type CompressParams struct {
	Format   string         `json:"format,omitempty"`
	Q        int            `json:"q,omitempty"`
	Strip    bool           `json:"strip,omitempty"`
	Extra    map[string]any `json:"-"`
	Delivery *Delivery      `json:"delivery,omitempty"`
}
