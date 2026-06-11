package report

import (
	"detector/internal/inspection/domain/inspector"
	"detector/internal/report/domain"
)

type Processor struct {
}

func NewReportProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) ProcessUserReport(domain.Report) error {

	return nil
}

func (p *Processor) ProcessInspectorReport(domain.Report, inspector.InspectionResult) error {

	return nil
}
