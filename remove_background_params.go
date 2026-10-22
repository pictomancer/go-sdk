package pictomancer

// RemoveBackgroundParams tunes the remove_background operation. Zero values
// are omitted; the server defaults Format to png. Format may be png, webp or
// avif (jpeg has no alpha). Q applies to webp and avif only.
type RemoveBackgroundParams struct {
	Format   string         `json:"format,omitempty"`
	Q        int            `json:"q,omitempty"`
	Extra    map[string]any `json:"-"`
	Delivery *Delivery      `json:"delivery,omitempty"`
}
