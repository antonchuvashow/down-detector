package inspectionport

import (
	inspectiondto "detector/internal/inspection/application/dto"
	routedomain "detector/internal/route/domain"
)

type RouteMethodRepo interface {
	SaveRouteMethod(routeMethod inspectiondto.RouteMethod) error
	GetRouteMethod(routeID routedomain.RouteID) (*inspectiondto.RouteMethod, error)
}
