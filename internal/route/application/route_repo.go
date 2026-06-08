package routeapplication

import routedomain "detector/internal/route/domain"

type RouteRepository interface {
	GetAllRoutes() ([]routedomain.Route, error)
	Get(id routedomain.RouteID) (routedomain.Route, error)
	Add(route routedomain.Route) error
	Update(route routedomain.Route) error
	Delete(id routedomain.RouteID) error
}
