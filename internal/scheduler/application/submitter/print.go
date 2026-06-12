package submitter

import (
	"go.uber.org/zap"

	"detector/internal/inspector/domain"
	"detector/internal/route/domain"
)

type PrintSubmitter struct {
	logger *zap.Logger
}

func NewPrintSubmitter(logger *zap.Logger) *PrintSubmitter {
	return &PrintSubmitter{
		logger: logger,
	}
}

func (p *PrintSubmitter) Submit(result inspector.Result, _ route.ID) error {
	p.logger.Info("Inspection result", zap.Any("result", result))
	return nil
}
