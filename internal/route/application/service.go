package routeapp

import (
	"fmt"

	"detector/internal/route/application/dto"
	"detector/internal/route/domain"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllRoutes() ([]route.Route, error) {
	return s.repo.GetAllRoutes()
}

func (s *Service) Get(id route.ID) (route.Route, error) {
	return s.repo.Get(id)
}

func (s *Service) Add(routeCommand routedto.AddCommand) (route.Route, error) {
	id, err := newRouteID()
	if err != nil {
		return route.Route{}, fmt.Errorf("could not generate RouteID: %w", err)
	}

	r := route.Route{
		ID:  id,
		URL: routeCommand.URL,
	}

	return r, s.repo.Add(r)
}

func (s *Service) Update(route route.Route) error {
	return s.repo.Update(route)
}

func (s *Service) Delete(id route.ID) error {
	return s.repo.Delete(id)
}

func newRouteID() (route.ID, error) {
	routeID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return route.ID(routeID.String()), nil
}
