package apiroute

import (
	"detector/internal/route/domain"
)

type Response struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func NewResponse(route route.Route) Response {
	return Response{
		ID:  string(route.ID),
		URL: route.URL.String(),
	}
}

type ListResponse struct {
	Routes []Response `json:"routes"`
	Total  int        `json:"total"`
}

func NewListResponse(routes []route.Route) ListResponse {
	responses := make([]Response, len(routes))
	for i, r := range routes {
		responses[i] = NewResponse(r)
	}
	return ListResponse{
		Routes: responses,
		Total:  len(responses),
	}
}
