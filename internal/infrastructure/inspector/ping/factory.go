package ping

import (
	"encoding/json"
	"fmt"
	"time"

	"detector/internal/inspection/domain/inspector"
)

type pingInspectorConfigJSON struct {
	PingCount *int           `json:"ping_count"`
	Interval  *time.Duration `json:"interval"`
	Timeout   *time.Duration `json:"timeout"`
	Threshold *float64       `json:"threshold"`
}

type PingInspectorFactory struct {
}

func (p *PingInspectorFactory) Create(_ *inspector.FactoryRegistry, config inspector.Config) (inspector.Inspector, error) {
	cfg, ok := config.(PingInspectorConfig)
	if !ok {
		return nil, fmt.Errorf("ping inspector factory: invalid config of type %T", config)
	}
	return NewInspector(cfg), nil

}

func (p *PingInspectorFactory) Marshal(_ *inspector.FactoryRegistry, object inspector.Inspector) (inspector.SerializedConfig, error) {
	// It may cause potential boilerplate and compatability issues.
	obj, ok := object.(*PingInspector)
	if !ok {
		return nil, fmt.Errorf("ping inspector factory: invalid object of type %T", object)
	}

	marshal, err := json.Marshal(pingInspectorConfigJSON(obj.config))
	if err != nil {
		return nil, fmt.Errorf("ping inspector factory: failed to marshal config: %w", err)
	}

	return marshal, nil
}

func (p *PingInspectorFactory) Unmarshal(_ *inspector.FactoryRegistry, data inspector.SerializedConfig) (inspector.Inspector, error) {
	var configJSON pingInspectorConfigJSON
	err := json.Unmarshal(data, &configJSON)
	if err != nil {
		return nil, fmt.Errorf("ping inspector factory: failed to unmarshal config: %w", err)
	}

	config := PingInspectorConfig(configJSON)
	return NewInspector(config), nil
}
