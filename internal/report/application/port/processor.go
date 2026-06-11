package port

import (
	"detector/internal/inspection/domain/inspector"
	"detector/internal/report/domain"
)

type ReportProcessor interface {
	ProcessUserReport(domain.Report) error
	ProcessInspectorReport(domain.Report, inspector.InspectionResult) error
}
