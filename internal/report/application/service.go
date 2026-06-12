package reportapp

import (
	"fmt"

	"detector/internal/inspector/domain"
	"detector/internal/report/domain"
)

type Service struct {
	processor Processor
	repo      Saver
}

func NewService(processor Processor, repo Saver) *Service {
	return &Service{
		processor: processor,
		repo:      repo,
	}
}

func (r *Service) SubmitUserReport(report report.Report) error {
	err := r.processor.ProcessUserReport(report)
	if err != nil {
		return fmt.Errorf("report service: %w", err)
	}
	err = r.repo.Save(report)
	if err != nil {
		return fmt.Errorf("report service: %w", err)
	}
	return nil
}

func (r *Service) SubmitInspectorReport(report report.Report, inspectorResult inspector.Result) error {
	err := r.processor.ProcessInspectorReport(report, inspectorResult)
	if err != nil {
		return fmt.Errorf("report service: %w", err)
	}
	err = r.repo.Save(report)
	if err != nil {
		return fmt.Errorf("report service: %w", err)
	}
	return nil
}
