package inspector

import (
	"time"
)

type InspectionStatus string
type ExtraInspectionInfo any

const (
	InspectionStatusSuccess InspectionStatus = "success"
	InspectionStatusError   InspectionStatus = "error"
)

type InspectionResult struct {
	Status InspectionStatus
	Start  time.Time
	End    time.Time
	Config Config
	Extra  ExtraInspectionInfo
}
