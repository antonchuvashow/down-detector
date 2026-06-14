package routeapp

import (
	"fmt"

	"detector/internal/route/application/dto"
	"detector/internal/route/domain"

	"github.com/google/uuid"
)

type Service struct {
	repo    Repository
	eventCh chan<- Event
}

func NewService(repo Repository, eventCh chan<- Event) *Service {
	return &Service{repo: repo, eventCh: eventCh}
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

	err = s.repo.Add(r)
	if err != nil {
		return route.Route{}, err
	}

	s.eventCh <- Event{Type: EventTypeCreate, Route: &r, RouteID: id}
	return r, nil
}

func (s *Service) Update(route route.Route) error {
	err := s.repo.Update(route)
	if err != nil {
		return err
	}

	s.eventCh <- Event{Type: EventTypeUpdate, Route: &route, RouteID: route.ID}
	return nil
}

func (s *Service) Delete(id route.ID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	s.eventCh <- Event{Type: EventTypeDelete, Route: nil, RouteID: id}
	return nil
}

func newRouteID() (route.ID, error) {
	routeID, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return route.ID(routeID.String()), nil
}
