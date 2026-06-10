package http

import (
	"net/http"
	"time"
)

type HttpExtraInspectionInfo struct {
	IsTimeout  bool
	StatusCode int
}

type HttpInspectorConfig struct {
	Timeout       *time.Duration
	ExpectedCodes map[int]struct{}
	Method        *string
	Header        http.Header
}

func NewInspectorConfig() *HttpInspectorConfig {
	expectedCodes := make(map[int]struct{})
	expectedCodes[200] = struct{}{}
	expectedCodes[204] = struct{}{}

	return &HttpInspectorConfig{
		Timeout:       new(time.Second),
		ExpectedCodes: expectedCodes,
		Method:        new("HEAD"),
		Header:        make(http.Header),
	}
}
