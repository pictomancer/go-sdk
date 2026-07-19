package pictomancer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sonirico/withttp"
)

const DefaultBaseURL = "https://api.pictomancer.ai"

// Client is a thin client for the Pictomancer.ai REST API.
type Client struct {
	adapter     withttp.Client
	endpoint    *withttp.Endpoint
	apiKey      string
	agentWallet string
	baseURL     string
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		adapter: withttp.NetHttp(),
		baseURL: DefaultBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	// Trailing slash matters: URI() joins with url.JoinPath, and a base
	// with an empty path yields a request line without the leading slash
	// ("GET v1/info"), which servers reject with 400.
	c.endpoint = withttp.NewEndpoint("pictomancer").
		Request(withttp.BaseURL(c.baseURL + "/"))
	return c
}

func (c *Client) Info(ctx context.Context) (InfoResponse, error) {
	res, err := c.roundTrip(ctx, http.MethodGet, "/v1/info", nil)
	if err != nil {
		return InfoResponse{}, err
	}
	return decodeJSON[InfoResponse](res.body)
}

func (c *Client) Usage(ctx context.Context) (UsageResponse, error) {
	res, err := c.roundTrip(ctx, http.MethodGet, "/v1/usage", nil)
	if err != nil {
		return UsageResponse{}, err
	}
	return decodeJSON[UsageResponse](res.body)
}

func (c *Client) Analyze(ctx context.Context, source string) (AnalyzeResponse, error) {
	res, err := c.roundTrip(ctx, http.MethodPost, "/v1/analyze", map[string]any{"source": source})
	if err != nil {
		return AnalyzeResponse{}, err
	}
	return decodeJSON[AnalyzeResponse](res.body)
}

func (c *Client) Resize(ctx context.Context, source string, params ResizeParams) (OpResult, error) {
	body, err := buildBody(source, params, params.Extra)
	if err != nil {
		return OpResult{}, err
	}
	return c.op(ctx, "/v1/resize", body, params.Delivery)
}

func (c *Client) Compress(ctx context.Context, source string, params CompressParams) (OpResult, error) {
	body, err := buildBody(source, params, params.Extra)
	if err != nil {
		return OpResult{}, err
	}
	return c.op(ctx, "/v1/compress", body, params.Delivery)
}

func (c *Client) Convert(ctx context.Context, source, format string, params ConvertParams) (OpResult, error) {
	body, err := buildBody(source, params, params.Extra)
	if err != nil {
		return OpResult{}, err
	}
	body["format"] = format
	return c.op(ctx, "/v1/convert", body, params.Delivery)
}

func (c *Client) Crop(ctx context.Context, source string, x, y, width, height int, params CropParams) (OpResult, error) {
	body, err := buildBody(source, params, params.Extra)
	if err != nil {
		return OpResult{}, err
	}
	body["x"] = x
	body["y"] = y
	body["width"] = width
	body["height"] = height
	return c.op(ctx, "/v1/crop", body, params.Delivery)
}

func (c *Client) Pipeline(ctx context.Context, source string, operations []PipelineOperation, delivery *Delivery) (OpResult, error) {
	body := map[string]any{"source": source, "operations": operations}
	if delivery != nil {
		body["delivery"] = delivery
	}
	return c.op(ctx, "/v1/pipeline", body, delivery)
}

// op runs an image operation. The requested delivery mode determines the
// result shape: inline responses are raw bytes, the rest are JSON receipts.
func (c *Client) op(ctx context.Context, path string, body map[string]any, delivery *Delivery) (OpResult, error) {
	res, err := c.roundTrip(ctx, http.MethodPost, path, body)
	if err != nil {
		return OpResult{}, err
	}
	if delivery.inline() {
		return OpResult{Bytes: res.body}, nil
	}
	receipt, err := decodeJSON[map[string]any](res.body)
	if err != nil {
		return OpResult{}, err
	}
	return OpResult{Receipt: receipt}, nil
}

type httpResult struct {
	status int
	body   []byte
}

func (c *Client) roundTrip(ctx context.Context, method, path string, payload map[string]any) (httpResult, error) {
	call := withttp.NewCall[struct{}](c.adapter).
		URI(path).
		Method(method).
		Header("User-Agent", userAgent(), true).
		ReadBody()
	if c.apiKey != "" {
		call.Header("Authorization", "Bearer "+c.apiKey, true)
	}
	if c.agentWallet != "" {
		call.Header("X-Agent-Wallet", c.agentWallet, true)
	}
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return httpResult{}, fmt.Errorf("marshal request body: %w", err)
		}
		call.ContentType("application/json").RawBody(raw)
	}

	if err := call.CallEndpoint(ctx, c.endpoint); err != nil {
		return httpResult{}, fmt.Errorf("call %s: %w", path, err)
	}
	status := call.Res.Status()
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		return httpResult{}, newAPIError(status, errorDetail(call.BodyRaw))
	}
	return httpResult{status: status, body: call.BodyRaw}, nil
}

// buildBody merges the params struct (via its JSON tags) with free-form
// extras, then sets the source. Delivery rides along through the params
// struct's own delivery tag.
func buildBody(source string, params any, extra map[string]any) (map[string]any, error) {
	body := map[string]any{}
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal params: %w", err)
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, fmt.Errorf("unmarshal params: %w", err)
	}
	for k, v := range extra {
		body[k] = v
	}
	body["source"] = source
	return body, nil
}

func decodeJSON[T any](body []byte) (T, error) {
	var out T
	if err := json.Unmarshal(body, &out); err != nil {
		return out, fmt.Errorf("decode response: %w", err)
	}
	return out, nil
}

func errorDetail(body []byte) string {
	var parsed struct {
		Detail string `json:"detail"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || parsed.Detail == "" {
		return string(body)
	}
	return parsed.Detail
}
