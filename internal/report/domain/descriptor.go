package domain

import "detector/internal/inspection/domain/inspector"

type Descriptor struct {
	Source SourceType
	Result inspector.InspectionResult
}
