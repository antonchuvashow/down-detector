package routeapplication

import (
	routedto "detector/internal/route/application/dto"
	routedomain "detector/internal/route/domain"
	"fmt"

	"github.com/google/uuid"
)

type RouteService struct {
	repo RouteRepository
}

func NewRouteService(repo RouteRepository) *RouteService {
	return &RouteService{repo: repo}
}

func (s *RouteService) GetAllRoutes() ([]routedomain.Route, error) {
	return s.repo.GetAllRoutes()
}

func (s *RouteService) Get(id routedomain.RouteID) (routedomain.Route, error) {
	return s.repo.Get(id)
}

func (s *RouteService) Add(routeCommand routedto.AddRouteCommand) (routedomain.Route, error) {
	id, err := newRouteID()
	if err != nil {
		return routedomain.Route{}, fmt.Errorf("could not generate RouteID: %w", err)
	}

	route := routedomain.Route{
		ID:  id,
		URL: routeCommand.URL,
	}

	return route, s.repo.Add(route)
}

func (s *RouteService) Update(route routedomain.Route) error {
	return s.repo.Update(route)
}

func (s *RouteService) Delete(id routedomain.RouteID) error {
	return s.repo.Delete(id)
}

func newRouteID() (routedomain.RouteID, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return routedomain.RouteID(""), err
	}

	return routedomain.RouteID(uuid.String()), nil
}
