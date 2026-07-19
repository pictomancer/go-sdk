package pictomancer

// FormatOption describes one tunable knob of an output format.
type FormatOption struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	DefaultStr  string `json:"default_str"`
	Description string `json:"description"`
	Min         *int   `json:"min,omitempty"`
	Max         *int   `json:"max,omitempty"`
}
