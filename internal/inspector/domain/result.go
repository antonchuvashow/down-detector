package inspector

import (
	"time"
)

type ResultStatus string
type ResultExtraInfo any

const (
	ResultStatusSuccess ResultStatus = "success"
	ResultStatusError   ResultStatus = "error"
)

type Result struct {
	Status ResultStatus
	Start  time.Time
	End    time.Time
	Config Config
	Extra  ResultExtraInfo
}
