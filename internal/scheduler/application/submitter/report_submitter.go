package submitter

import (
	"fmt"
	"time"

	"detector/internal/infrastructure/inspector/composite"
	"detector/internal/infrastructure/inspector/http"
	"detector/internal/infrastructure/inspector/ping"
	"detector/internal/inspection/domain/inspector"
	"detector/internal/report/application/service"
	"detector/internal/report/domain"
	routedomain "detector/internal/route/domain"
)

type ReportSubmitter struct {
	reportService *service.ReportService
	descriptor    domain.Descriptor
}

func NewReportSubmitter(reportService *service.ReportService, descriptor domain.Descriptor) *ReportSubmitter {
	return &ReportSubmitter{
		reportService: reportService,
		descriptor:    descriptor,
	}
}

func (s *ReportSubmitter) Submit(result inspector.InspectionResult, routeID routedomain.RouteID) error {
	collectedErrors := make(map[domain.ErrorType]struct{})
	var latency time.Duration
	collectErrorsAndLatency(result, collectedErrors, &latency)

	report := domain.Report{
		Success:    result.Status == inspector.InspectionStatusSuccess,
		RouteID:    routeID,
		ErrorTypes: collectedErrors,
		Time:       result.Start,
		Descriptor: s.descriptor,
		Summary: domain.Summary{
			Latency: latency,
		},
	}

	err := s.reportService.SubmitUserReport(report)
	if err != nil {
		return fmt.Errorf("submit user report: %w", err)
	}
	return nil
}

func collectErrorsAndLatency(inspectionResult inspector.InspectionResult, collectedErrors map[domain.ErrorType]struct{}, latency *time.Duration) {
	// For composite inspector latency defined as the sum of latencies of internal inspectors
	switch v := inspectionResult.Extra.(type) {
	case ping.ExtraInspectionInfo:
		if inspectionResult.Status == inspector.InspectionStatusError {
			collectedErrors[domain.ErrorTypeNetwork] = struct{}{}
		}
		*latency += v.AvgRtt // needed for composite inspector
	case http.ExtraInspectionInfo:
		if inspectionResult.Status == inspector.InspectionStatusError {
			if v.IsTimeout {
				collectedErrors[domain.ErrorTypeNetwork] = struct{}{}
			} else if v.StatusCode >= 500 {
				collectedErrors[domain.ErrorTypeServer] = struct{}{}
				collectedErrors[domain.ErrorTypeWebAccess] = struct{}{}
				collectedErrors[domain.ErrorTypeMobileAccess] = struct{}{}
			} else {
				collectedErrors[domain.ErrorTypeUnknown] = struct{}{}
			}
		}
		*latency += inspectionResult.End.Sub(inspectionResult.Start) // needed for composite inspector
	case composite.ExtraInspectionInfo:
		for _, result := range v.Results {
			collectErrorsAndLatency(result, collectedErrors, latency)
		}
	default:
		collectedErrors[domain.ErrorTypeUnknown] = struct{}{}
	}
}
