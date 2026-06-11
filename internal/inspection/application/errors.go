package application

import (
	"fmt"

	routedomain "detector/internal/route/domain"
)

type ErrInspectorNotFound struct {
	RouteID routedomain.RouteID
}

func (e *ErrInspectorNotFound) Error() string {
	return fmt.Sprintf("inspector not found for route ID: %s", string(e.RouteID))
}
