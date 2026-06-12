package reportapp

import "detector/internal/report/domain"

type Saver interface {
	Save(report report.Report) error
}
