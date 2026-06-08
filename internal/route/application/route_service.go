package application

import (
	"detector/internal/route/application/dto"
	"detector/internal/route/domain"
	"fmt"

	"github.com/google/uuid"
)

type RouteService struct {
	repo RouteRepository
}

func NewRouteService(repo RouteRepository) *RouteService {
	return &RouteService{repo: repo}
}

func (s *RouteService) GetAllRoutes() ([]domain.Route, error) {
	return s.repo.GetAllRoutes()
}

func (s *RouteService) Get(id domain.RouteID) (domain.Route, error) {
	return s.repo.Get(id)
}

func (s *RouteService) Add(routeCommand dto.AddRouteCommand) (domain.Route, error) {
	id, err := newRouteID()
	if err != nil {
		return domain.Route{}, fmt.Errorf("could not generate RouteID: %w", err)
	}

	route := domain.Route{
		ID:  id,
		URL: routeCommand.URL,
	}

	return route, s.repo.Add(route)
}

func (s *RouteService) Update(route domain.Route) error {
	return s.repo.Update(route)
}

func (s *RouteService) Delete(id domain.RouteID) error {
	return s.repo.Delete(id)
}

func newRouteID() (domain.RouteID, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return domain.RouteID(""), err
	}

	return domain.RouteID(uuid.String()), nil
}
