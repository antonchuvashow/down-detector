package schedulerapp

import (
	"detector/internal/inspector/domain"
	"detector/internal/route/domain"
)

type InspectorResultSubmitter interface {
	Submit(result inspector.Result, routeID route.ID) error
}
