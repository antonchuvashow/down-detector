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
)

type ReportSubmitter struct {
	reportService *service.ReportService
}

func NewReportSubmitter(reportService *service.ReportService) *ReportSubmitter {
	return &ReportSubmitter{reportService: reportService}
}

func (s *ReportSubmitter) Submit(inspectionResult inspector.InspectionResult) error {
	collectedErrors := make(map[domain.ErrorType]struct{})
	var latency time.Duration
	collectErrorsAndLatency(inspectionResult, collectedErrors, &latency)

	report := domain.Report{
		Success:    inspectionResult.Status == inspector.InspectionStatusSuccess,
		ErrorTypes: collectedErrors,
		Time:       inspectionResult.Start,
		Descriptor: domain.Descriptor{Source: domain.SourceTypeInspector},
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
