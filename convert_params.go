package pictomancer

// ConvertParams tunes the convert operation. Zero values are omitted.
// Effort is a pointer because 0 is a valid AVIF encoder value distinct
// from unset (API default 2); use Int to build it.
type ConvertParams struct {
	Q        int            `json:"q,omitempty"`
	Strip    bool           `json:"strip,omitempty"`
	Lossless bool           `json:"lossless,omitempty"`
	Effort   *int           `json:"effort,omitempty"`
	Extra    map[string]any `json:"-"`
	Delivery *Delivery      `json:"delivery,omitempty"`
}
