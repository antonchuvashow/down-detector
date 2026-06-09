package inspector

import (
	routedomain "detector/internal/route/domain"
)

type Inspector interface {
	Inspect(routedomain.Route) (InspectionResult, error)
}
