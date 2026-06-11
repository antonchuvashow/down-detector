package inspectiondto

import (
	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"
)

// RouteMethod is a link between route and inspector snapshot
// described as a FactoryKey and a SerializedConfig.
type RouteMethod struct {
	RouteID          routedomain.RouteID
	FactoryKey       inspector.FactoryKey
	SerializedConfig inspector.SerializedConfig
}
