package routeapp

import "detector/internal/route/domain"

type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeUpdate EventType = "update"
	EventTypeDelete EventType = "delete"
)

type Event struct {
	Type    EventType
	RouteID route.ID
	Route   *route.Route
}
