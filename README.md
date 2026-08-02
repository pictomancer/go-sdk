# pictomancer go-sdk

Go SDK for [Pictomancer.ai](https://pictomancer.ai) - a thin client for the REST API at `https://api.pictomancer.ai`, built on [withttp](https://github.com/sonirico/withttp).

## Install

```bash
go get github.com/pictomancer/go-sdk
```

## Usage

```go
package main

import (
	"context"
	"os"

	pictomancer "github.com/pictomancer/go-sdk"
)

func main() {
	client := pictomancer.NewClient(
		pictomancer.WithAPIKey("your-api-key"),
	)
	ctx := context.Background()

	info, err := client.Info(ctx)
	if err != nil {
		panic(err)
	}
	_ = info

	meta, err := client.Analyze(ctx, "https://example.com/image.jpg")
	if err != nil {
		panic(err)
	}
	_ = meta.SizeBytes

	result, err := client.Compress(ctx, "https://example.com/image.jpg", pictomancer.CompressParams{
		Format: "webp",
		Q:      80,
		Strip:  true,
	})
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("out.webp", result.Bytes, 0o644); err != nil {
		panic(err)
	}
}
```

Sources can be an image URL, a base64 string, or a `data:` URI. For local files, readers, or in-memory bytes:

```go
source, err := pictomancer.SourceFromPath("photo.jpg")
source, err = pictomancer.SourceFromReader(file)
source = pictomancer.SourceFromBytes(data)
```

Timeouts and cancellation flow through `context.Context`.

### Operations

```go
client.Info(ctx)
client.Usage(ctx)
client.Analyze(ctx, source)
client.Resize(ctx, source, pictomancer.ResizeParams{Scale: 0.5, Format: "webp"})
client.Compress(ctx, source, pictomancer.CompressParams{Q: 80})
client.Convert(ctx, source, "avif", pictomancer.ConvertParams{Q: 50, Effort: pictomancer.Int(2)})
client.Crop(ctx, source, 0, 0, 100, 100, pictomancer.CropParams{Format: "png"})
client.Pipeline(ctx, source, []pictomancer.PipelineOperation{
	{Type: "resize", Params: map[string]string{"scale": "0.5"}},
	{Type: "convert", Params: map[string]string{"format": "webp"}},
}, nil)
```

Operations return an `OpResult`: `Bytes` holds the optimized image for inline delivery, `Receipt` holds the JSON receipt for `put_url`/`callback` deliveries.

### Perceptual quality target

Instead of guessing a `Q`, ask compress/convert for the smallest file with SSIM >= target. The server binary-searches the encoder quality and reports the outcome in `X-Pictomancer-Quality-*` headers, surfaced as `OpResult.Quality`:

```go
result, err := client.Compress(ctx, source, pictomancer.CompressParams{
	Format:        "webp", // required with QualityTarget
	QualityTarget: 0.95,   // 0 < v <= 1; mutually exclusive with Q
})
if err != nil {
	panic(err)
}
if result.Quality != nil {
	// result.Quality.Achieved (e.g. 0.9530), .QFinal, .Encodes
}
```

Supported for `jpeg`, `webp` and `avif` outputs; on convert it is also invalid with `Lossless: true`. Not available inside pipelines. `Quality` is nil when no search ran - either no `QualityTarget` was sent, or the input already met the target and came back untouched (`X-Pig-Billed: 0`).

### Delivery targets

```go
// Presigned PUT (S3/R2/GCS/Azure). No cloud credentials reach Pictomancer.
result, err := client.Compress(ctx, source, pictomancer.CompressParams{
	Format:   "webp",
	Delivery: pictomancer.NewPutURLDelivery(presignedURL),
})

// POST to your endpoint, HMAC-signed (X-Pig-Signature: sha256=<hex>).
result, err = client.Convert(ctx, source, "avif", pictomancer.ConvertParams{
	Delivery: pictomancer.NewCallbackDelivery(
		"https://hooks.example.com/pig?token=...",
		pictomancer.WithDeliverySecret(secret),
	),
})
```

### Options

```go
pictomancer.WithAPIKey("...")        // Bearer token
pictomancer.WithBaseURL("...")       // defaults to https://api.pictomancer.ai
pictomancer.WithAgentWallet("0x...") // X-Agent-Wallet for x402 tracking
pictomancer.WithAdapter(withttp.Fasthttp()) // swap the HTTP backend
```

### Errors

Non-2xx responses return `*pictomancer.APIError`:

```go
var apiErr *pictomancer.APIError
if errors.As(err, &apiErr) && apiErr.Status == 402 {
	// free tier exhausted: pay per request (x402) or use an API key
}
```

## Development

```bash
go test ./...
```

## License

MIT
