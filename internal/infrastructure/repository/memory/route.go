package repository

import (
	routeapplication "detector/internal/route/application"
	routedomain "detector/internal/route/domain"
)

type MemoryRouteRepository struct {
	routes map[routedomain.RouteID]routedomain.Route
}

func NewMemoryRouteRepository() *MemoryRouteRepository {
	return &MemoryRouteRepository{
		routes: make(map[routedomain.RouteID]routedomain.Route),
	}
}

func (r *MemoryRouteRepository) GetAllRoutes() ([]routedomain.Route, error) {
	result := make([]routedomain.Route, 0, len(r.routes))
	for _, route := range r.routes {
		result = append(result, route)
	}
	return result, nil
}

func (r *MemoryRouteRepository) Get(id routedomain.RouteID) (routedomain.Route, error) {
	route, exists := r.routes[id]
	if !exists {
		return routedomain.Route{}, routeapplication.ErrRouteNotFound{ID: id}
	}
	return route, nil
}

func (r *MemoryRouteRepository) Add(route routedomain.Route) error {
	if _, exists := r.routes[route.ID]; exists {
		return routeapplication.ErrRouteAlreadyExists{ID: route.ID}
	}
	r.routes[route.ID] = route
	return nil
}

func (r *MemoryRouteRepository) Update(route routedomain.Route) error {
	if _, exists := r.routes[route.ID]; !exists {
		return routeapplication.ErrRouteNotFound{ID: route.ID}
	}
	r.routes[route.ID] = route
	return nil
}

func (r *MemoryRouteRepository) Delete(id routedomain.RouteID) error {
	if _, exists := r.routes[id]; !exists {
		return routeapplication.ErrRouteNotFound{ID: id}
	}
	delete(r.routes, id)
	return nil
}
