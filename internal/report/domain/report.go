package report

import (
	"time"

	"detector/internal/route/domain"
)

type Report struct {
	Success    bool
	RouteID    route.ID
	ErrorTypes map[ErrorType]struct{}
	Time       time.Time
	Descriptor Descriptor
	Summary    Summary
}
