package inspectorapp

import (
	"fmt"

	"detector/internal/route/domain"
)

type ErrNotFound struct {
	RouteID route.ID
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("inspector not found for route ID: %s", string(e.RouteID))
}
