package application

import "detector/internal/inspection/domain/inspector"

type InspectionResultSubmitter interface {
	Submit(result inspector.InspectionResult) error
}
