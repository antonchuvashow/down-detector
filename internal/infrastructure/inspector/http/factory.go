package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"detector/internal/inspection/domain/inspector"
)

type httpInspectorConfigJSON struct {
	Timeout       *time.Duration `json:"timeout"`
	ExpectedCodes []int          `json:"expected_codes"`
	Method        *string        `json:"method"`
	Header        http.Header    `json:"header"`
}

type HttpInspectorFactory struct {
}

func (h *HttpInspectorFactory) Create(_ *inspector.FactoryRegistry, config inspector.Config) (inspector.Inspector, error) {
	cfg, ok := config.(HttpInspectorConfig)
	if !ok {
		return nil, fmt.Errorf("http inspector factory: invalid config of type %T", config)
	}

	return NewInspector(cfg), nil
}

func (h *HttpInspectorFactory) Marshal(_ *inspector.FactoryRegistry, object inspector.Inspector) (inspector.SerializedConfig, error) {
	obj, ok := object.(*HttpInspector)
	if !ok {
		return nil, fmt.Errorf("http inspector factory: invalid object of type %T", object)
	}

	configJSON := httpInspectorConfigJSON{
		Timeout:       obj.config.Timeout,
		Method:        obj.config.Method,
		Header:        obj.config.Header,
		ExpectedCodes: setToSlice(obj.config.ExpectedCodes),
	}

	marshal, err := json.Marshal(configJSON)
	if err != nil {
		return nil, fmt.Errorf("http inspector factory: failed to marshal config: %w", err)
	}

	return marshal, nil
}

func (h *HttpInspectorFactory) Unmarshal(_ *inspector.FactoryRegistry, data inspector.SerializedConfig) (inspector.Inspector, error) {
	var configJSON httpInspectorConfigJSON
	err := json.Unmarshal(data, &configJSON)
	if err != nil {
		return nil, fmt.Errorf("http inspector factory: failed to unmarshal config: %w", err)
	}

	config := HttpInspectorConfig{
		Timeout:       configJSON.Timeout,
		Method:        configJSON.Method,
		Header:        configJSON.Header,
		ExpectedCodes: sliceToSet(configJSON.ExpectedCodes),
	}

	return NewInspector(config), nil
}

func sliceToSet(values []int) map[int]struct{} {
	set := make(map[int]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func setToSlice(values map[int]struct{}) []int {
	set := make([]int, 0, len(values))
	for value := range values {
		set = append(set, value)
	}
	return set
}
