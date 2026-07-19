package pictomancer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeliveryConstructors(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		got  *Delivery
		want *Delivery
	}

	cases := []testCase{
		{
			name: "inline builds the default target",
			got:  NewInlineDelivery(),
			want: &Delivery{Mode: "inline"},
		},
		{
			name: "put url builds a put_url target",
			got:  NewPutURLDelivery("https://bucket.example.com/key?sig=1"),
			want: &Delivery{Mode: "put_url", PutURL: "https://bucket.example.com/key?sig=1"},
		},
		{
			name: "put url includes signed headers",
			got: NewPutURLDelivery(
				"https://bucket.example.com/key",
				WithDeliveryHeaders(map[string]string{"x-amz-acl": "private"}),
			),
			want: &Delivery{
				Mode:    "put_url",
				PutURL:  "https://bucket.example.com/key",
				Headers: map[string]string{"x-amz-acl": "private"},
			},
		},
		{
			name: "callback builds a callback_url target",
			got:  NewCallbackDelivery("https://hooks.example.com/pig?token=t"),
			want: &Delivery{Mode: "callback_url", CallbackURL: "https://hooks.example.com/pig?token=t"},
		},
		{
			name: "callback includes hmac secret",
			got:  NewCallbackDelivery("https://hooks.example.com/pig", WithDeliverySecret("hmac-me")),
			want: &Delivery{
				Mode:        "callback_url",
				CallbackURL: "https://hooks.example.com/pig",
				Secret:      "hmac-me",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, tc.got)
		})
	}
}
