package pictomancer

// Int returns a pointer to v, for optional int params such as
// ConvertParams.Effort.
func Int(v int) *int {
	return &v
}
