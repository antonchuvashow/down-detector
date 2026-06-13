package reportapp

import (
	"go.uber.org/zap"

	"detector/internal/inspector/domain"
	"detector/internal/report/domain"
)

type Processor interface {
	ProcessUserReport(report.Report) error
	ProcessInspectorReport(report.Report, inspector.Result) error
}

type PrintProcessor struct {
	logger *zap.Logger
}

func NewPrintProcessor(logger *zap.Logger) *PrintProcessor {
	return &PrintProcessor{logger: logger}
}

func (p *PrintProcessor) ProcessUserReport(report report.Report) error {
	p.logger.Info("ProcessUserReport", zap.Bool("success", report.Success))
	return nil
}

func (p *PrintProcessor) ProcessInspectorReport(report report.Report, _ inspector.Result) error {
	p.logger.Info("ProcessInspectorReport", zap.Bool("success", report.Success))
	return nil
}
