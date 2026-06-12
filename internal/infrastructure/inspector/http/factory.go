package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"detector/internal/inspector/domain"
)

type inspectorConfigJSON struct {
	Timeout       *time.Duration `json:"timeout"`
	ExpectedCodes []int          `json:"expected_codes"`
	Method        *string        `json:"method"`
	Header        http.Header    `json:"header"`
}

type InspectorFactory struct {
}

func (h *InspectorFactory) Marshal(_ *inspector.FactoryRegistry, object inspector.Inspector) (inspector.SerializedConfig, error) {
	obj, ok := object.(*Inspector)
	if !ok {
		return nil, fmt.Errorf("http inspector factory: invalid object of type %T", object)
	}

	configJSON := inspectorConfigJSON{
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

func (h *InspectorFactory) Unmarshal(_ *inspector.FactoryRegistry, data inspector.SerializedConfig) (inspector.Inspector, error) {
	var configJSON inspectorConfigJSON
	err := json.Unmarshal(data, &configJSON)
	if err != nil {
		return nil, fmt.Errorf("http inspector factory: failed to unmarshal config: %w", err)
	}

	config := InspectorConfig{
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
