package port

import "detector/internal/report/domain"

type ReportSaver interface {
	Save(report domain.Report)
}
