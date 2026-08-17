package pictomancer

import (
	"fmt"
	"runtime"
	"strings"
)

const Version = "0.6.0"

func userAgent() string {
	return fmt.Sprintf(
		"pictomancer-go/%s go/%s",
		Version,
		strings.TrimPrefix(runtime.Version(), "go"),
	)
}
