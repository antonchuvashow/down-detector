package application

import (
	"go.uber.org/zap"

	"detector/internal/inspection/domain/inspector"
	"detector/internal/report/domain"
)

type PrintProcessor struct {
	logger *zap.Logger
}

func NewReportProcessor(logger *zap.Logger) *PrintProcessor {
	return &PrintProcessor{logger: logger}
}

func (p *PrintProcessor) ProcessUserReport(report domain.Report) error {
	p.logger.Info("ProcessUserReport", zap.Any("report", report))
	return nil
}

func (p *PrintProcessor) ProcessInspectorReport(report domain.Report, object inspector.InspectionResult) error {
	p.logger.Info("ProcessInspectorReport", zap.Any("report", report))
	return nil
}
