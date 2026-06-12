package inmemory

import (
	"detector/internal/route/application"
	"detector/internal/route/domain"
)

type RouteRepository struct {
	routes map[route.ID]route.Route
}

func NewRouteRepository() *RouteRepository {
	return &RouteRepository{
		routes: make(map[route.ID]route.Route),
	}
}

func (rr *RouteRepository) GetAllRoutes() ([]route.Route, error) {
	result := make([]route.Route, 0, len(rr.routes))
	for _, rt := range rr.routes {
		result = append(result, rt)
	}
	return result, nil
}

func (rr *RouteRepository) Get(id route.ID) (route.Route, error) {
	rt, exists := rr.routes[id]
	if !exists {
		return route.Route{}, routeapp.ErrNotFound{ID: id}
	}
	return rt, nil
}

func (rr *RouteRepository) Add(route route.Route) error {
	if _, exists := rr.routes[route.ID]; exists {
		return routeapp.ErrAlreadyExists{ID: route.ID}
	}
	rr.routes[route.ID] = route
	return nil
}

func (rr *RouteRepository) Update(route route.Route) error {
	if _, exists := rr.routes[route.ID]; !exists {
		return routeapp.ErrNotFound{ID: route.ID}
	}
	rr.routes[route.ID] = route
	return nil
}

func (rr *RouteRepository) Delete(id route.ID) error {
	if _, exists := rr.routes[id]; !exists {
		return routeapp.ErrNotFound{ID: id}
	}
	delete(rr.routes, id)
	return nil
}
