package http

import (
	"net/http"
	"time"
)

type ExtraInspectionInfo struct {
	IsTimeout  bool
	StatusCode int
}

type InspectorConfig struct {
	Timeout       *time.Duration
	ExpectedCodes map[int]struct{}
	Method        *string
	Header        http.Header
}

func NewInspectorConfig() *InspectorConfig {
	expectedCodes := make(map[int]struct{})
	expectedCodes[200] = struct{}{}
	expectedCodes[204] = struct{}{}

	return &InspectorConfig{
		Timeout:       new(time.Second),
		ExpectedCodes: expectedCodes,
		Method:        new("HEAD"),
		Header:        make(http.Header),
	}
}
