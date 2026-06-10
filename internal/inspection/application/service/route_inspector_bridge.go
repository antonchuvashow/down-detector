package inspectionservice

import (
	"fmt"

	inspectiondto "detector/internal/inspection/application/dto"
	inspectionport "detector/internal/inspection/application/port"
	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"
)

type RouteInspectorBridge struct {
	factoryRegistry *inspector.FactoryRegistry
	repo            inspectionport.RouteMethodRepo
}

func NewRouteInspectorBridge(factoryRegistry *inspector.FactoryRegistry, repository inspectionport.RouteMethodRepo) RouteInspectorBridge {
	return RouteInspectorBridge{
		factoryRegistry: factoryRegistry,
		repo:            repository,
	}
}

func (b *RouteInspectorBridge) Register(routeID routedomain.RouteID, key inspector.FactoryKey, config inspector.Config) error {
	// TODO: Check if route exists, if not - return error

	factory, err := b.factoryRegistry.Get(key)
	if err != nil {
		return fmt.Errorf("route inspector bridge: failed to get inspector factory: %w", err)
	}

	object, err := factory.Create(b.factoryRegistry, config)
	if err != nil {
		return fmt.Errorf("route inspector bridge: failed to create inspector: %w", err)
	}

	serializedConfig, err := factory.Marshal(b.factoryRegistry, object)
	if err != nil {
		return fmt.Errorf("route inspector bridge: failed to marshal inspector config: %w", err)
	}

	routeMethod := inspectiondto.RouteMethod{
		RouteID:          routeID,
		FactoryKey:       key,
		SerializedConfig: serializedConfig,
	}

	err = b.repo.SaveRouteMethod(routeMethod)
	if err != nil {
		return fmt.Errorf("route inspector bridge: failed to save route method: %w", err)
	}

	return nil
}

func (b *RouteInspectorBridge) FindInspector(routeID routedomain.RouteID) (inspector.Inspector, error) {
	routeMethod, err := b.repo.GetRouteMethod(routeID)
	if err != nil {
		return nil, fmt.Errorf("route inspector bridge: failed to get route method: %w", err)
	}

	factory, err := b.factoryRegistry.Get(routeMethod.FactoryKey)
	if err != nil {
		return nil, fmt.Errorf("route inspector bridge: failed to get inspector factory: %w", err)
	}

	object, err := factory.Unmarshal(b.factoryRegistry, routeMethod.SerializedConfig)
	if err != nil {
		return nil, fmt.Errorf("route inspector bridge: failed to unmarshal inspector config: %w", err)
	}

	return object, nil
}
