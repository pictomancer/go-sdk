package pictomancer

// ResizeParams tunes the resize operation. Zero values are omitted.
// Use Scale for uniform scaling or ScaleX/ScaleY for independent axes.
type ResizeParams struct {
	Scale    float64        `json:"scale,omitempty"`
	ScaleX   float64        `json:"scale_x,omitempty"`
	ScaleY   float64        `json:"scale_y,omitempty"`
	Format   string         `json:"format,omitempty"`
	Extra    map[string]any `json:"-"`
	Delivery *Delivery      `json:"delivery,omitempty"`
}
