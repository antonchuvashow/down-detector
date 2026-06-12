package inspectordto

import (
	"detector/internal/inspector/domain"
	"detector/internal/route/domain"
)

// RouteAssignment is a link between route and inspector snapshot
// described as a FactoryKey and a SerializedConfig.
type RouteAssignment struct {
	RouteID          route.ID
	FactoryKey       inspector.FactoryKey
	SerializedConfig inspector.SerializedConfig
}
