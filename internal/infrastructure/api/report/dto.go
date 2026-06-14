package apireport

import (
	"errors"

	"detector/internal/report/domain"
)

// SubmitRequest is the JSON body for POST /reports.
type SubmitRequest struct {
	// RouteID is the UUID of the monitored route being reported.
	RouteID string `json:"route_id" binding:"required"`

	// Success indicates whether the user considers the service reachable.
	Success bool `json:"success"`

	// ErrorTypes is a list of error type strings chosen by the user.
	// Valid values are defined in report.ErrorType.
	ErrorTypes []string `json:"error_types"`

	// Source identifies whether this is a user or automated report.
	// Defaults to "user" when omitted.
	Source string `json:"source"`

	// Platform is the device/OS type reported by the client.
	// Allowed: unknown, linux, windows, android, ios.
	Platform string `json:"platform"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	LatencyMs float64 `json:"latency_ms"`
}

// Validate performs domain-level validation of the request.
func (r *SubmitRequest) Validate() error {
	if r.RouteID == "" {
		return errors.New("route_id is required")
	}
	if !r.Success && len(r.ErrorTypes) == 0 {
		return errors.New("at least one error_type must be provided when success is false")
	}
	return nil
}

// ErrorTypeSet converts the slice of error type strings into the domain set.
func (r *SubmitRequest) ErrorTypeSet() map[report.ErrorType]struct{} {
	set := make(map[report.ErrorType]struct{}, len(r.ErrorTypes))
	for _, et := range r.ErrorTypes {
		set[report.ErrorType(et)] = struct{}{}
	}
	return set
}

// ErrorResponse is a generic JSON error envelope.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ValidationErrorResponse wraps a validation error.
func ValidationErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}
