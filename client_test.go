package pictomancer

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const testSource = "https://example.com/image.jpg"

var testImageBytes = []byte{0xff, 0xd8, 0xff, 0xe0}

type recordedRequest struct {
	method  string
	path    string
	headers http.Header
	body    map[string]any
}

type testResponse struct {
	status      int
	contentType string
	body        []byte
}

type testFixture struct {
	client   *Client
	requests *[]recordedRequest
}

func newTestFixture(t *testing.T, response testResponse, opts ...Option) testFixture {
	t.Helper()

	requests := &[]recordedRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		recorded := recordedRequest{method: r.Method, path: r.URL.Path, headers: r.Header.Clone()}
		if len(raw) > 0 {
			require.NoError(t, json.Unmarshal(raw, &recorded.body))
		}
		*requests = append(*requests, recorded)

		contentType := response.contentType
		if contentType == "" {
			contentType = "application/json"
		}
		w.Header().Set("content-type", contentType)
		status := response.status
		if status == 0 {
			status = http.StatusOK
		}
		w.WriteHeader(status)
		_, err = w.Write(response.body)
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := NewClient(append([]Option{WithBaseURL(server.URL)}, opts...)...)
	return testFixture{client: client, requests: requests}
}

func jsonBody(t *testing.T, v any) []byte {
	t.Helper()
	raw, err := json.Marshal(v)
	require.NoError(t, err)
	return raw
}

func TestClientHeaders(t *testing.T) {
	t.Parallel()

	t.Run("sends bearer auth when api key is set", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{body: []byte(`{"formats":[]}`)}, WithAPIKey("sk-123"))

		_, err := fx.client.Info(context.Background())

		require.NoError(t, err)
		require.Equal(t, "Bearer sk-123", (*fx.requests)[0].headers.Get("Authorization"))
	})

	t.Run("omits auth header without api key", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{body: []byte(`{"formats":[]}`)})

		_, err := fx.client.Info(context.Background())

		require.NoError(t, err)
		require.Empty(t, (*fx.requests)[0].headers.Get("Authorization"))
	})

	t.Run("sends user agent telemetry", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{body: []byte(`{"formats":[]}`)})

		_, err := fx.client.Info(context.Background())

		require.NoError(t, err)
		require.Regexp(t, `^pictomancer-go/\d+\.\d+\.\d+ go/`, (*fx.requests)[0].headers.Get("User-Agent"))
	})

	t.Run("sends agent wallet when configured", func(t *testing.T) {
		t.Parallel()
		wallet := "0xabababababababababababababababababababab"
		fx := newTestFixture(t, testResponse{body: []byte(`{"formats":[]}`)}, WithAgentWallet(wallet))

		_, err := fx.client.Info(context.Background())

		require.NoError(t, err)
		require.Equal(t, wallet, (*fx.requests)[0].headers.Get("X-Agent-Wallet"))
	})
}

func TestClientRequestBodies(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		act      func(c *Client) error
		response testResponse
		wantPath string
		wantBody map[string]any
	}

	cases := []testCase{
		{
			name: "analyze posts source only",
			act: func(c *Client) error {
				_, err := c.Analyze(context.Background(), testSource)
				return err
			},
			response: testResponse{body: []byte(`{"size_bytes":1}`)},
			wantPath: "/v1/analyze",
			wantBody: map[string]any{"source": testSource},
		},
		{
			name: "resize sends scale params",
			act: func(c *Client) error {
				_, err := c.Resize(context.Background(), testSource, ResizeParams{Scale: 0.5, Format: "webp"})
				return err
			},
			wantPath: "/v1/resize",
			wantBody: map[string]any{"source": testSource, "scale": 0.5, "format": "webp"},
		},
		{
			name: "resize sends independent axes",
			act: func(c *Client) error {
				_, err := c.Resize(context.Background(), testSource, ResizeParams{ScaleX: 0.5, ScaleY: 2.0})
				return err
			},
			wantPath: "/v1/resize",
			wantBody: map[string]any{"source": testSource, "scale_x": 0.5, "scale_y": 2.0},
		},
		{
			name: "compress sends quality params",
			act: func(c *Client) error {
				_, err := c.Compress(context.Background(), testSource, CompressParams{Format: "jpeg", Q: 60, Strip: true})
				return err
			},
			wantPath: "/v1/compress",
			wantBody: map[string]any{"source": testSource, "format": "jpeg", "q": float64(60), "strip": true},
		},
		{
			name: "convert sends format and effort zero",
			act: func(c *Client) error {
				_, err := c.Convert(context.Background(), testSource, "avif", ConvertParams{Q: 50, Effort: Int(0)})
				return err
			},
			wantPath: "/v1/convert",
			wantBody: map[string]any{"source": testSource, "format": "avif", "q": float64(50), "effort": float64(0)},
		},
		{
			name: "crop sends region including zero origin",
			act: func(c *Client) error {
				_, err := c.Crop(context.Background(), testSource, 0, 0, 100, 50, CropParams{Format: "png"})
				return err
			},
			wantPath: "/v1/crop",
			wantBody: map[string]any{
				"source": testSource, "x": float64(0), "y": float64(0),
				"width": float64(100), "height": float64(50), "format": "png",
			},
		},
		{
			name: "pipeline sends the operation chain",
			act: func(c *Client) error {
				ops := []PipelineOperation{
					{Type: "resize", Params: map[string]string{"scale": "0.5"}},
					{Type: "convert", Params: map[string]string{"format": "webp"}},
				}
				_, err := c.Pipeline(context.Background(), testSource, ops, nil)
				return err
			},
			wantPath: "/v1/pipeline",
			wantBody: map[string]any{
				"source": testSource,
				"operations": []any{
					map[string]any{"type": "resize", "params": map[string]any{"scale": "0.5"}},
					map[string]any{"type": "convert", "params": map[string]any{"format": "webp"}},
				},
			},
		},
		{
			name: "extra params pass through",
			act: func(c *Client) error {
				_, err := c.Compress(context.Background(), testSource, CompressParams{
					Format: "png",
					Extra:  map[string]any{"palette": true},
				})
				return err
			},
			wantPath: "/v1/compress",
			wantBody: map[string]any{"source": testSource, "format": "png", "palette": true},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			response := tc.response
			if response.body == nil {
				response = testResponse{contentType: "image/jpeg", body: testImageBytes}
			}
			fx := newTestFixture(t, response)

			err := tc.act(fx.client)

			require.NoError(t, err)
			require.Len(t, *fx.requests, 1)
			got := (*fx.requests)[0]
			require.Equal(t, http.MethodPost, got.method)
			require.Equal(t, tc.wantPath, got.path)
			require.Equal(t, tc.wantBody, got.body)
		})
	}
}

func TestClientDelivery(t *testing.T) {
	t.Parallel()

	t.Run("attaches put_url target to the body", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{body: []byte(`{"sha256":"abc"}`)})
		delivery := NewPutURLDelivery(
			"https://bucket.example.com/key?sig=1",
			WithDeliveryHeaders(map[string]string{"x-amz-acl": "private"}),
		)

		_, err := fx.client.Resize(context.Background(), testSource, ResizeParams{Scale: 0.5, Delivery: delivery})

		require.NoError(t, err)
		require.Equal(t, map[string]any{
			"mode":    "put_url",
			"put_url": "https://bucket.example.com/key?sig=1",
			"headers": map[string]any{"x-amz-acl": "private"},
		}, (*fx.requests)[0].body["delivery"])
	})

	t.Run("attaches callback target with secret", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{body: []byte(`{"status":200}`)})
		delivery := NewCallbackDelivery("https://hooks.example.com/pig?token=t", WithDeliverySecret("hmac-me"))

		_, err := fx.client.Compress(context.Background(), testSource, CompressParams{Delivery: delivery})

		require.NoError(t, err)
		require.Equal(t, map[string]any{
			"mode":         "callback_url",
			"callback_url": "https://hooks.example.com/pig?token=t",
			"secret":       "hmac-me",
		}, (*fx.requests)[0].body["delivery"])
	})

	t.Run("inline delivery returns bytes", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{contentType: "image/jpeg", body: testImageBytes})

		result, err := fx.client.Resize(context.Background(), testSource, ResizeParams{Scale: 0.5})

		require.NoError(t, err)
		require.Equal(t, testImageBytes, result.Bytes)
		require.Nil(t, result.Receipt)
	})

	t.Run("put_url delivery returns a receipt", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{body: jsonBody(t, map[string]any{"sha256": "abc", "bytes_written": 42})})
		delivery := NewPutURLDelivery("https://bucket.example.com/key")

		result, err := fx.client.Resize(context.Background(), testSource, ResizeParams{Scale: 0.5, Delivery: delivery})

		require.NoError(t, err)
		require.Nil(t, result.Bytes)
		require.Equal(t, map[string]any{"sha256": "abc", "bytes_written": float64(42)}, result.Receipt)
	})
}

func TestClientResponses(t *testing.T) {
	t.Parallel()

	t.Run("info decodes format specs", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{
			body: jsonBody(t, map[string]any{"formats": []map[string]any{{"id": "avif", "suffix": ".avif"}}}),
		})

		info, err := fx.client.Info(context.Background())

		require.NoError(t, err)
		require.Len(t, info.Formats, 1)
		require.Equal(t, "avif", info.Formats[0].ID)
	})

	t.Run("usage decodes free tier state", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{
			body: jsonBody(t, map[string]any{"identity": "ip:1.2.3.4", "requests_used": 3, "free_remaining": 47, "is_free": true}),
		})

		usage, err := fx.client.Usage(context.Background())

		require.NoError(t, err)
		require.Equal(t, "ip:1.2.3.4", usage.Identity)
		require.Equal(t, 3, usage.RequestsUsed)
		require.True(t, usage.IsFree)
	})

	t.Run("analyze decodes size", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{body: []byte(`{"size_bytes":1024}`)})

		meta, err := fx.client.Analyze(context.Background(), testSource)

		require.NoError(t, err)
		require.Equal(t, int64(1024), meta.SizeBytes)
	})
}

func TestClientErrors(t *testing.T) {
	t.Parallel()

	t.Run("non 2xx returns APIError with detail", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{status: http.StatusPaymentRequired, body: []byte(`{"detail":"payment required"}`)})

		_, err := fx.client.Resize(context.Background(), testSource, ResizeParams{Scale: 0.5})

		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusPaymentRequired, apiErr.Status)
		require.Equal(t, "payment required", apiErr.Detail)
	})

	t.Run("non json error body falls back to raw text", func(t *testing.T) {
		t.Parallel()
		fx := newTestFixture(t, testResponse{status: http.StatusBadGateway, contentType: "text/plain", body: []byte("boom")})

		_, err := fx.client.Info(context.Background())

		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusBadGateway, apiErr.Status)
		require.Equal(t, "boom", apiErr.Detail)
	})
}
