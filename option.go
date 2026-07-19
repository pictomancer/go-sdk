package pictomancer

import (
	"strings"

	"github.com/sonirico/withttp"
)

type Option func(*Client)

// WithAPIKey sets the Bearer token (Authorization: Bearer ...).
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithBaseURL overrides the default https://api.pictomancer.ai.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(url, "/")
	}
}

// WithAgentWallet sets the X-Agent-Wallet identity for x402 tracking.
func WithAgentWallet(wallet string) Option {
	return func(c *Client) {
		c.agentWallet = wallet
	}
}

// WithAdapter swaps the HTTP backend (e.g. withttp.Fasthttp()).
// Defaults to net/http.
func WithAdapter(adapter withttp.Client) Option {
	return func(c *Client) {
		c.adapter = adapter
	}
}
