package pictomancer

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	headerQualityTarget   = "X-Pictomancer-Quality-Target"
	headerQualityAchieved = "X-Pictomancer-Quality-Achieved"
	headerQualityQFinal   = "X-Pictomancer-Quality-Q-Final"
	headerQualityEncodes  = "X-Pictomancer-Quality-Encodes"
)

// QualityReport is the outcome of a quality_target SSIM search, parsed from
// the X-Pictomancer-Quality-* response headers. Absent (nil on OpResult)
// when no search ran: no quality_target requested, or the input already met
// the target and came back untouched (X-Pig-Billed: 0).
type QualityReport struct {
	Target   float64
	Achieved float64
	QFinal   int
	Encodes  int
}

// newQualityReport builds the report from response headers. Returns nil
// when the headers are absent.
func newQualityReport(header http.Header) (*QualityReport, error) {
	if header.Get(headerQualityTarget) == "" {
		return nil, nil
	}
	target, err := strconv.ParseFloat(header.Get(headerQualityTarget), 64)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", headerQualityTarget, err)
	}
	achieved, err := strconv.ParseFloat(header.Get(headerQualityAchieved), 64)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", headerQualityAchieved, err)
	}
	qFinal, err := strconv.Atoi(header.Get(headerQualityQFinal))
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", headerQualityQFinal, err)
	}
	encodes, err := strconv.Atoi(header.Get(headerQualityEncodes))
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", headerQualityEncodes, err)
	}
	return &QualityReport{Target: target, Achieved: achieved, QFinal: qFinal, Encodes: encodes}, nil
}
