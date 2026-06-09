package inspector

import (
	"time"
)

type InspectionStatus string
type ExtraInspectionInfo any

const (
	StatusSuccess InspectionStatus = "success"
	StatusError   InspectionStatus = "error"
)

type InspectionResult struct {
	Status InspectionStatus
	Start  time.Time
	End    time.Time
	Config Config
	Extra  ExtraInspectionInfo
}
