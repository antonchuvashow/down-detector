package inspectiondto

import (
	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"
)

type RouteMethod struct {
	RouteID          routedomain.RouteID
	FactoryKey       inspector.FactoryKey
	SerializedConfig inspector.SerializedConfig
}
