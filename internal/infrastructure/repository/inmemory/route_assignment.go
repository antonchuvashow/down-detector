package inmemory

import (
	"detector/internal/inspector/application"
	"detector/internal/inspector/application/dto"
	"detector/internal/route/domain"
)

type RouteAssignmentRepository struct {
	routeAssignments map[route.ID]inspectordto.RouteAssignment
}

func NewRouteAssignmentRepository() *RouteAssignmentRepository {
	return &RouteAssignmentRepository{
		routeAssignments: make(map[route.ID]inspectordto.RouteAssignment),
	}
}

func (r *RouteAssignmentRepository) Save(routeAssignment inspectordto.RouteAssignment) error {
	r.routeAssignments[routeAssignment.RouteID] = routeAssignment
	return nil
}

func (r *RouteAssignmentRepository) Get(routeID route.ID) (*inspectordto.RouteAssignment, error) {
	routeAssignment, exists := r.routeAssignments[routeID]
	if !exists {
		return nil, &inspectorapp.ErrNotFound{RouteID: routeID}
	}
	return &routeAssignment, nil
}

func (r *RouteAssignmentRepository) Delete(routeID route.ID) error {
	_, err := r.Get(routeID)
	if err != nil {
		return err
	}

	delete(r.routeAssignments, routeID)
	return nil
}
