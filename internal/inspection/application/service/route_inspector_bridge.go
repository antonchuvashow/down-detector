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

func (b *RouteInspectorBridge) Register(routeID routedomain.RouteID, instance inspector.Inspector) error {
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

	routeMethod := inspectiondto.RouteMethod{
		RouteID:          routeID,
		FactoryKey:       factoryKey,
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
