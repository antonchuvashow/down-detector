package inmemory

import "detector/internal/report/domain"

type ReportSaver struct {
	reports []report.Report
}

func NewReportSaver() *ReportSaver {
	return &ReportSaver{}
}

func (r *ReportSaver) Save(report report.Report) error {
	r.reports = append(r.reports, report)
	return nil
}
