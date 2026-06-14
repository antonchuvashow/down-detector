package apireport

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"detector/internal/report/application"
	"detector/internal/report/domain"
	routedomain "detector/internal/route/domain"
)

type Handler struct {
	service *reportapp.Service
}

func NewHandler(service *reportapp.Service) *Handler {
	return &Handler{service: service}
}

// RegisterReports registers all user-report-related endpoints.
func (h *Handler) RegisterReports(r *gin.RouterGroup) {
	reports := r.Group("/reports")
	reports.POST("", h.Submit)
}

// Submit handles POST /reports – a user-submitted report.
func (h *Handler) Submit(c *gin.Context) {
	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse(err))
		return
	}

	platform := report.PlatformType(req.Platform)
	if platform == "" || platform == report.PlatformTypeUnknown {
		platform = detectPlatformFromUA(c.GetHeader("User-Agent"))
	}

	r := report.Report{
		Success:    req.Success,
		RouteID:    routedomain.ID(req.RouteID),
		ErrorTypes: req.ErrorTypeSet(),
		Time:       time.Now().UTC(),
		Descriptor: report.Descriptor{
			Source:    report.SourceTypeUser,
			Platform:  platform,
			IP:        net.ParseIP(c.ClientIP()),
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
		Summary: report.Summary{
			Latency: time.Duration(req.LatencyMs) * time.Millisecond,
		},
	}

	if err := h.service.SubmitUserReport(r); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func detectPlatformFromUA(ua string) report.PlatformType {
	ua = strings.ToLower(ua)

	switch {
	case strings.Contains(ua, "android"):
		return report.PlatformTypeAndroid
	case strings.Contains(ua, "iphone"),
		strings.Contains(ua, "ipad"),
		strings.Contains(ua, "ipod"):
		return report.PlatformTypeIOS
	case strings.Contains(ua, "windows"):
		return report.PlatformTypeWindows
	case strings.Contains(ua, "linux") && !strings.Contains(ua, "android"):
		return report.PlatformTypeLinux
	default:
		return report.PlatformTypeUnknown
	}
}
