package repository

import "detector/internal/report/domain"

type ReportSaver struct {
	reports []domain.Report
}

func NewReportSaver() *ReportSaver {
	return &ReportSaver{}
}

func (r *ReportSaver) Save(report domain.Report) error {
	r.reports = append(r.reports, report)
	return nil
}
