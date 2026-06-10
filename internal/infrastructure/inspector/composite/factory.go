package composite

import (
	"encoding/json"
	"fmt"

	"detector/internal/inspection/domain/inspector"
)

type compositeInspectorConfigJSON struct {
	Inspectors map[string]factoryConfigJSON `json:"inspectors"`
}

type factoryConfigJSON struct {
	FactoryKey string `json:"factory_key"`
	Config     string `json:"config"`
}

type CompositeInspectorFactory struct {
}

func (c *CompositeInspectorFactory) Create(registry *inspector.FactoryRegistry, config inspector.Config) (inspector.Inspector, error) {
	cfg, ok := config.(CompositeInspectorConfig)
	if !ok {
		return nil, fmt.Errorf("composite inspector factory: invalid config of type %T", config)
	}

	return NewCompositeInspector(cfg), nil
}

func (c *CompositeInspectorFactory) Marshal(registry *inspector.FactoryRegistry, object inspector.Inspector) (inspector.SerializedConfig, error) {
	obj, ok := object.(*CompositeInspector)
	if !ok {
		return nil, fmt.Errorf("composite inspector factory: invalid object of type %T", object)
	}
	configJSON := compositeInspectorConfigJSON{
		Inspectors: make(map[string]factoryConfigJSON),
	}

	for name, instance := range obj.config.Inspectors {
		factoryKey, err := registry.Find(instance)
		if err != nil {
			return nil, fmt.Errorf("composite inspector factory: could not find factory key for type %T", instance)
		}
		factory, err := registry.Get(factoryKey)
		if err != nil {
			return nil, fmt.Errorf("composite inspector factory: could not find factory for key %s", factoryKey)
		}
		marshal, err := factory.Marshal(registry, instance)
		if err != nil {
			return nil, fmt.Errorf("composite inspector factory: could not marshal inspector of type %T", instance)
		}
		configJSON.Inspectors[name] = factoryConfigJSON{
			FactoryKey: string(factoryKey),
			Config:     string(marshal),
		}
	}

	marshal, err := json.Marshal(configJSON)
	if err != nil {
		return nil, fmt.Errorf("ping inspector factory: failed to marshal config: %w", err)
	}

	return marshal, nil
}

func (c *CompositeInspectorFactory) Unmarshal(registry *inspector.FactoryRegistry, data inspector.SerializedConfig) (inspector.Inspector, error) {
	var configJSON compositeInspectorConfigJSON
	err := json.Unmarshal(data, &configJSON)
	if err != nil {
		return nil, fmt.Errorf("composite inspector factory: failed to unmarshal config: %w", err)
	}
	config := CompositeInspectorConfig{
		Inspectors: make(map[string]inspector.Inspector),
	}

	for name, fcJSON := range configJSON.Inspectors {
		factory, err := registry.Get(inspector.FactoryKey(fcJSON.FactoryKey))
		if err != nil {
			return nil, fmt.Errorf("composite inspector factory: could not find factory for key %s: %w", fcJSON.FactoryKey, err)
		}
		instance, err := factory.Unmarshal(registry, []byte(fcJSON.Config))
		if err != nil {
			return nil, fmt.Errorf("composite inspector factory: could not unmarshal config: %w", err)
		}
		config.Inspectors[name] = instance
	}

	return NewCompositeInspector(config), nil
}
