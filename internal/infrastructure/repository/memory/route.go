package repository

import (
	"detector/internal/route/application"
	"detector/internal/route/domain"
)

type MemoryRouteRepository struct {
	routes map[domain.RouteID]domain.Route
}

func NewMemoryRouteRepository() *MemoryRouteRepository {
	return &MemoryRouteRepository{
		routes: make(map[domain.RouteID]domain.Route),
	}
}

func (r *MemoryRouteRepository) GetAllRoutes() ([]domain.Route, error) {
	result := make([]domain.Route, 0, len(r.routes))
	for _, route := range r.routes {
		result = append(result, route)
	}
	return result, nil
}

func (r *MemoryRouteRepository) Get(id domain.RouteID) (domain.Route, error) {
	route, exists := r.routes[id]
	if !exists {
		return domain.Route{}, application.ErrRouteNotFound{ID: id}
	}
	return route, nil
}

func (r *MemoryRouteRepository) Add(route domain.Route) error {
	if _, exists := r.routes[route.ID]; exists {
		return application.ErrRouteAlreadyExists{ID: route.ID}
	}
	r.routes[route.ID] = route
	return nil
}

func (r *MemoryRouteRepository) Update(route domain.Route) error {
	if _, exists := r.routes[route.ID]; !exists {
		return application.ErrRouteNotFound{ID: route.ID}
	}
	r.routes[route.ID] = route
	return nil
}

func (r *MemoryRouteRepository) Delete(id domain.RouteID) error {
	if _, exists := r.routes[id]; !exists {
		return application.ErrRouteNotFound{ID: id}
	}
	delete(r.routes, id)
	return nil
}
