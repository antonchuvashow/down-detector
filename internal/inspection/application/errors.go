package application

import (
	routedomain "detector/internal/route/domain"
	"fmt"
)

type ErrInspectorNotFound struct {
	RouteID routedomain.RouteID
}

func (e *ErrInspectorNotFound) Error() string {
	return fmt.Sprintf("inspector not found for route ID: %s", string(e.RouteID))
}
