package application

import (
	"detector/internal/route/domain"
	"fmt"
)

type ErrRouteNotFound struct {
	ID domain.RouteID
}

func (e ErrRouteNotFound) Error() string {
	return fmt.Sprintf("route with ID %v not found", e.ID)
}

type ErrRouteAlreadyExists struct {
	ID domain.RouteID
}

func (e ErrRouteAlreadyExists) Error() string {
	return fmt.Sprintf("route with ID %v already exists", e.ID)
}
