package application

import "detector/internal/route/domain"

type RouteRepository interface {
	GetAllRoutes() ([]domain.Route, error)
	Get(id domain.RouteID) (domain.Route, error)
	Add(route domain.Route) error
	Update(route domain.Route) error
	Delete(id domain.RouteID) error
}
