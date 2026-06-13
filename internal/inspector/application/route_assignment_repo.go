package inspectorapp

import (
	"detector/internal/inspector/application/dto"
	"detector/internal/route/domain"
)

type RouteAssignmentRepository interface {
	Save(routeAssignment inspectordto.RouteAssignment) error
	Get(routeID route.ID) (*inspectordto.RouteAssignment, error)
	Delete(routeID route.ID) error
}
