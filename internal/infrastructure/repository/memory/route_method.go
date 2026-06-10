package repository

import (
	"detector/internal/inspection/application"
	inspectiondto "detector/internal/inspection/application/dto"
	routedomain "detector/internal/route/domain"
)

type MemoryRouteMethodRepository struct {
	routeMethods map[routedomain.RouteID]inspectiondto.RouteMethod
}

func NewMemoryRouteMethodRepository() *MemoryRouteMethodRepository {
	return &MemoryRouteMethodRepository{
		routeMethods: make(map[routedomain.RouteID]inspectiondto.RouteMethod),
	}
}

func (r *MemoryRouteMethodRepository) SaveRouteMethod(routeMethod inspectiondto.RouteMethod) error {
	r.routeMethods[routeMethod.RouteID] = routeMethod
	return nil
}

func (r *MemoryRouteMethodRepository) GetRouteMethod(routeID routedomain.RouteID) (*inspectiondto.RouteMethod, error) {
	routeMethod, exists := r.routeMethods[routeID]
	if !exists {
		return nil, &application.ErrInspectorNotFound{RouteID: routeID}
	}
	return &routeMethod, nil
}
