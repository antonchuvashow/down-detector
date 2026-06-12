package routeapp

import "detector/internal/route/domain"

type Repository interface {
	GetAllRoutes() ([]route.Route, error)
	Get(id route.ID) (route.Route, error)
	Add(route route.Route) error
	Update(route route.Route) error
	Delete(id route.ID) error
}
