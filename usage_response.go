package pictomancer

// UsageResponse reports request usage and free tier status for the
// caller's identity (wallet header or IP fallback).
type UsageResponse struct {
	Identity      string `json:"identity"`
	RequestsUsed  int    `json:"requests_used"`
	FreeTierLimit int    `json:"free_tier_limit"`
	FreeRemaining int    `json:"free_remaining"`
	IsFree        bool   `json:"is_free"`
}
