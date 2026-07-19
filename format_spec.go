package pictomancer

// FormatSpec describes one output format (id, file suffix, options).
type FormatSpec struct {
	ID      string         `json:"id"`
	Suffix  string         `json:"suffix"`
	Options []FormatOption `json:"options"`
}
