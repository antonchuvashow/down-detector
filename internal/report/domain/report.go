package domain

import (
	"time"

	routedomain "detector/internal/route/domain"
)

type Report struct {
	Success    bool
	RouteID    routedomain.RouteID
	ErrorTypes map[ErrorType]struct{}
	Time       time.Time
	Descriptor Descriptor
	Summary    Summary
}
