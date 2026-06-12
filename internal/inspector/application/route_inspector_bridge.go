package inspectorapp

import (
	"fmt"

	"detector/internal/inspector/application/dto"
	"detector/internal/inspector/domain"
	"detector/internal/route/domain"
)

type RouteInspectorBridge struct {
	factoryRegistry *inspector.FactoryRegistry
	repo            RouteAssignmentRepository
}

func NewRouteInspectorBridge(factoryRegistry *inspector.FactoryRegistry, repository RouteAssignmentRepository) RouteInspectorBridge {
	return RouteInspectorBridge{
		factoryRegistry: factoryRegistry,
		repo:            repository,
	}
}

func (b *RouteInspectorBridge) Register(routeID route.ID, instance inspector.Inspector) error {
	// TODO: Check if route exists, if not - return error

	factoryKey, err := b.factoryRegistry.Find(instance)
	if err != nil {
		return fmt.Errorf("route inspector bridge: could not find factory key for inspector of type %T: %w", instance, err)
	}
	factory, err := b.factoryRegistry.Get(factoryKey)
	if err != nil {
		return fmt.Errorf("route inspector bridge: could not find factory for key %s: %w", factoryKey, err)
	}

	serializedConfig, err := factory.Marshal(b.factoryRegistry, instance)
	if err != nil {
		return fmt.Errorf("route inspector bridge: failed to marshal inspector config: %w", err)
	}

	assignment := inspectordto.RouteAssignment{
		RouteID:          routeID,
		FactoryKey:       factoryKey,
		SerializedConfig: serializedConfig,
	}

	err = b.repo.Save(assignment)
	if err != nil {
		return fmt.Errorf("route inspector bridge: failed to save route assignment: %w", err)
	}

	return nil
}

func (b *RouteInspectorBridge) FindInspector(routeID route.ID) (inspector.Inspector, error) {
	assignment, err := b.repo.Get(routeID)
	if err != nil {
		return nil, fmt.Errorf("route inspector bridge: failed to get route assignment: %w", err)
	}

	factory, err := b.factoryRegistry.Get(assignment.FactoryKey)
	if err != nil {
		return nil, fmt.Errorf("route inspector bridge: failed to get inspector factory: %w", err)
	}

	object, err := factory.Unmarshal(b.factoryRegistry, assignment.SerializedConfig)
	if err != nil {
		return nil, fmt.Errorf("route inspector bridge: failed to unmarshal inspector config: %w", err)
	}

	return object, nil
}
