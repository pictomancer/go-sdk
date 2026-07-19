package pictomancer

// InfoResponse lists the supported output formats and their options.
type InfoResponse struct {
	Formats []FormatSpec `json:"formats"`
}
