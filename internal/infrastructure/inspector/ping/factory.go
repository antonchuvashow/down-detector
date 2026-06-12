package ping

import (
	"encoding/json"
	"fmt"
	"time"

	"detector/internal/inspector/domain"
)

type inspectorConfigJSON struct {
	PingCount *int           `json:"ping_count"`
	Interval  *time.Duration `json:"interval"`
	Timeout   *time.Duration `json:"timeout"`
	Threshold *float64       `json:"threshold"`
}

type InspectorFactory struct {
}

func (p *InspectorFactory) Marshal(_ *inspector.FactoryRegistry, object inspector.Inspector) (inspector.SerializedConfig, error) {
	// It may cause potential boilerplate and compatability issues.
	obj, ok := object.(*Inspector)
	if !ok {
		return nil, fmt.Errorf("ping inspector factory: invalid object of type %T", object)
	}

	marshal, err := json.Marshal(inspectorConfigJSON(obj.config))
	if err != nil {
		return nil, fmt.Errorf("ping inspector factory: failed to marshal config: %w", err)
	}

	return marshal, nil
}

func (p *InspectorFactory) Unmarshal(_ *inspector.FactoryRegistry, data inspector.SerializedConfig) (inspector.Inspector, error) {
	var configJSON inspectorConfigJSON
	err := json.Unmarshal(data, &configJSON)
	if err != nil {
		return nil, fmt.Errorf("ping inspector factory: failed to unmarshal config: %w", err)
	}

	config := InspectorConfig(configJSON)
	return NewInspector(config), nil
}
