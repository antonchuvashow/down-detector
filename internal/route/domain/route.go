package routedomain

import "net/url"

type RouteID string

type Route struct {
	ID  RouteID
	URL url.URL
}
