package pictomancer

// CropParams tunes the crop operation. Three mutually exclusive modes:
// manual (X+Y+Width+Height), smart (Gravity + Width+Height, X/Y nil), trim
// (Trim=true, optional Threshold, X/Y/Width/Height nil). Autorot is valid
// in all three. X/Y/Width/Height are pointers because 0 is a legitimate
// corner distinct from unset; use Int to build them.
type CropParams struct {
	X         *int           `json:"x,omitempty"`
	Y         *int           `json:"y,omitempty"`
	Width     *int           `json:"width,omitempty"`
	Height    *int           `json:"height,omitempty"`
	Format    string         `json:"format,omitempty"`
	Gravity   string         `json:"gravity,omitempty"`
	Trim      bool           `json:"trim,omitempty"`
	Threshold float64        `json:"threshold,omitempty"`
	Autorot   bool           `json:"autorot,omitempty"`
	Extra     map[string]any `json:"-"`
	Delivery  *Delivery      `json:"delivery,omitempty"`
}
