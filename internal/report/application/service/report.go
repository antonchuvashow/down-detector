package service

import (
	"fmt"

	"detector/internal/inspection/domain/inspector"
	"detector/internal/report/application/port"
	"detector/internal/report/domain"
)

type ReportService struct {
	processor port.ReportProcessor
	repo      port.ReportSaver
}

func NewReportService(notifier port.ReportProcessor, repo port.ReportSaver) *ReportService {
	return &ReportService{
		processor: notifier,
		repo:      repo,
	}
}

func (r *ReportService) SubmitUserReport(report domain.Report) error {
	err := r.processor.ProcessUserReport(report)
	if err != nil {
		return fmt.Errorf("report service: %w", err)
	}
	r.repo.Save(report)
	return nil
}

func (r *ReportService) SubmitInspectorReport(report domain.Report, inspectorResult inspector.InspectionResult) error {
	err := r.processor.ProcessInspectorReport(report, inspectorResult)
	if err != nil {
		return fmt.Errorf("report service: %w", err)
	}
	r.repo.Save(report)
	return nil
}
