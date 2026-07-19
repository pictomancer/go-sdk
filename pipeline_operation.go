package pictomancer

// PipelineOperation is one step of a multi-op pipeline. Params values
// are strings, matching the API contract (e.g. {"scale": "0.5"}).
type PipelineOperation struct {
	Type   string            `json:"type"`
	Params map[string]string `json:"params"`
}
