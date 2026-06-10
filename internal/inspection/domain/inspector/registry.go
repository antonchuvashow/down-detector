package inspector

import (
	"fmt"
	"reflect"
)

type FactoryKey string

type FactoryRegistry struct {
	typeToKey map[reflect.Type]FactoryKey
	factories map[FactoryKey]InspectorFactory
}

func NewFactoryRegistry() FactoryRegistry {
	return FactoryRegistry{
		factories: make(map[FactoryKey]InspectorFactory),
		typeToKey: make(map[reflect.Type]FactoryKey),
	}
}

func (r *FactoryRegistry) Register(key FactoryKey, inspectorType reflect.Type, factory InspectorFactory) {
	if _, exists := r.factories[key]; exists {
		panic(fmt.Sprintf("factory with key %s already registered", string(key)))
	}

	r.factories[key] = factory
	r.typeToKey[inspectorType] = key
}

func (r *FactoryRegistry) Get(key FactoryKey) (InspectorFactory, error) {
	factory, exists := r.factories[key]
	if !exists {
		return nil, fmt.Errorf("factory with key %s not found", string(key))
	}
	return factory, nil
}

func (r *FactoryRegistry) Find(inspector Inspector) (FactoryKey, error) {
	inspectorType := reflect.TypeOf(inspector)
	key, ok := r.typeToKey[inspectorType]
	if !ok {
		return "", fmt.Errorf("could not determine factory key for type: %T", inspector)
	}

	return key, nil
}
