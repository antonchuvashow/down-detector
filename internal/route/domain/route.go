package route

import "net/url"

type ID string

type Route struct {
	ID  ID
	URL url.URL
}
