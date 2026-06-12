package submitter

import (
	"fmt"
	"time"

	"detector/internal/infrastructure/inspector/composite"
	"detector/internal/infrastructure/inspector/http"
	"detector/internal/infrastructure/inspector/ping"
	"detector/internal/inspector/domain"
	"detector/internal/report/application"
	"detector/internal/report/domain"
	"detector/internal/route/domain"
)

type ReportSubmitter struct {
	reportService *reportapp.Service
	descriptor    report.Descriptor
}

func NewReportSubmitter(reportService *reportapp.Service, descriptor report.Descriptor) *ReportSubmitter {
	return &ReportSubmitter{
		reportService: reportService,
		descriptor:    descriptor,
	}
}

func (s *ReportSubmitter) Submit(result inspector.Result, routeID route.ID) error {
	collectedErrors := make(map[report.ErrorType]struct{})
	var latency time.Duration
	collectErrorsAndLatency(result, collectedErrors, &latency)

	r := report.Report{
		Success:    result.Status == inspector.ResultStatusSuccess,
		RouteID:    routeID,
		ErrorTypes: collectedErrors,
		Time:       result.Start,
		Descriptor: s.descriptor,
		Summary: report.Summary{
			Latency: latency,
		},
	}

	err := s.reportService.SubmitUserReport(r)
	if err != nil {
		return fmt.Errorf("submit user r: %w", err)
	}
	return nil
}

func collectErrorsAndLatency(inspectionResult inspector.Result, collectedErrors map[report.ErrorType]struct{}, latency *time.Duration) {
	// For composite inspector latency defined as the sum of latencies of internal inspectors
	switch v := inspectionResult.Extra.(type) {
	case ping.ExtraInspectionInfo:
		if inspectionResult.Status == inspector.ResultStatusError {
			collectedErrors[report.ErrorTypeNetwork] = struct{}{}
		}
		*latency += v.AvgRtt // needed for composite inspector
	case http.ExtraInspectionInfo:
		if inspectionResult.Status == inspector.ResultStatusError {
			if v.IsTimeout {
				collectedErrors[report.ErrorTypeNetwork] = struct{}{}
			} else if v.StatusCode >= 500 {
				collectedErrors[report.ErrorTypeServer] = struct{}{}
				collectedErrors[report.ErrorTypeWebAccess] = struct{}{}
				collectedErrors[report.ErrorTypeMobileAccess] = struct{}{}
			} else {
				collectedErrors[report.ErrorTypeUnknown] = struct{}{}
			}
		}
		*latency += inspectionResult.End.Sub(inspectionResult.Start) // needed for composite inspector
	case composite.ExtraInspectionInfo:
		for _, result := range v.Results {
			collectErrorsAndLatency(result, collectedErrors, latency)
		}
	default:
		collectedErrors[report.ErrorTypeUnknown] = struct{}{}
	}
}
