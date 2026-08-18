package pictomancer

// ResizeParams tunes the resize operation. Zero values are omitted.
// Use Scale for uniform scaling or ScaleX/ScaleY for independent axes.
// Fill mode: set Width+Height instead, optionally with Gravity (one of
// "attention", "entropy", "centre"; defaults to "attention" server-side).
// Denoise (1-3), Equalize and Sharpen are opt-in enhancement modifiers.
type ResizeParams struct {
	Scale    float64        `json:"scale,omitempty"`
	ScaleX   float64        `json:"scale_x,omitempty"`
	ScaleY   float64        `json:"scale_y,omitempty"`
	Format   string         `json:"format,omitempty"`
	Width    *int           `json:"width,omitempty"`
	Height   *int           `json:"height,omitempty"`
	Gravity  string         `json:"gravity,omitempty"`
	Autorot  bool           `json:"autorot,omitempty"`
	Denoise  int            `json:"denoise,omitempty"`
	Equalize bool           `json:"equalize,omitempty"`
	Sharpen  bool           `json:"sharpen,omitempty"`
	Extra    map[string]any `json:"-"`
	Delivery *Delivery      `json:"delivery,omitempty"`
}
