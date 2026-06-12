package routeapp

import (
	"fmt"

	"detector/internal/route/domain"
)

type ErrNotFound struct {
	ID route.ID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("route with ID %v not found", e.ID)
}

type ErrAlreadyExists struct {
	ID route.ID
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("route with ID %v already exists", e.ID)
}
