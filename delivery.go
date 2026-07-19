package pictomancer

// Delivery selects where the optimized bytes go. Inline (the default)
// returns them in the response; put_url uploads to a customer-signed
// presigned PUT URL; callback_url POSTs them to a customer endpoint
// with an X-Pig-Sha256 integrity header and optional HMAC signature.
type Delivery struct {
	Mode        string            `json:"mode"`
	PutURL      string            `json:"put_url,omitempty"`
	CallbackURL string            `json:"callback_url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Secret      string            `json:"secret,omitempty"`
}

type DeliveryOption func(*Delivery)

// WithDeliveryHeaders forwards extra signed headers (x-amz-*, x-goog-*,
// x-ms-*) that the presigned URL or endpoint expects.
func WithDeliveryHeaders(headers map[string]string) DeliveryOption {
	return func(d *Delivery) {
		d.Headers = headers
	}
}

// WithDeliverySecret enables HMAC-SHA256 signing of callback bodies:
// the request carries X-Pig-Signature: sha256=<hex>, recomputable on the
// receiving end with the same secret. Used per request, never stored.
func WithDeliverySecret(secret string) DeliveryOption {
	return func(d *Delivery) {
		d.Secret = secret
	}
}

func NewInlineDelivery() *Delivery {
	return &Delivery{Mode: "inline"}
}

func NewPutURLDelivery(url string, opts ...DeliveryOption) *Delivery {
	d := &Delivery{Mode: "put_url", PutURL: url}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func NewCallbackDelivery(url string, opts ...DeliveryOption) *Delivery {
	d := &Delivery{Mode: "callback_url", CallbackURL: url}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *Delivery) inline() bool {
	return d == nil || d.Mode == "inline"
}
