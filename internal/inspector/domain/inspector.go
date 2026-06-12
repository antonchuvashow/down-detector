package inspector

import (
	"detector/internal/route/domain"
)

type Inspector interface {
	Inspect(route.Route) (Result, error)
}
