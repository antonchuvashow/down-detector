package apiroute

import (
	"fmt"
	"net/url"

	"detector/internal/route/application/dto"
	"detector/internal/route/domain"
)

type CreateRequest struct {
	URL string `json:"url" binding:"required"`
}

func (r CreateRequest) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("url is required")
	}
	if _, err := url.Parse(r.URL); err != nil {
		return fmt.Errorf("invalid url format")
	}
	return nil
}

func (r CreateRequest) ToCommand() routedto.AddCommand {
	u, _ := url.Parse(r.URL)
	return routedto.AddCommand{URL: *u}
}

type UpdateRequest struct {
	URL string `json:"url" binding:"required"`
}

func (r UpdateRequest) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("url is required")
	}
	if _, err := url.Parse(r.URL); err != nil {
		return fmt.Errorf("invalid url format")
	}
	return nil
}

func (r UpdateRequest) ToCommand(id string) route.Route {
	u, _ := url.Parse(r.URL)
	return route.Route{URL: *u, ID: route.ID(id)}
}
