package application

import (
	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"
)

type InspectionResultSubmitter interface {
	Submit(result inspector.InspectionResult, routeID routedomain.RouteID) error
}
