package apisuperset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"detector/internal/infrastructure/client/superset"
)

type Dashboard struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Handler struct {
	client     *superset.Client
	logger     *zap.Logger
	guest      superset.GuestDescriptor
	dashboards []Dashboard
}

func NewHandler(client *superset.Client, guestDescriptor superset.GuestDescriptor, dashboards []Dashboard, logger *zap.Logger) *Handler {
	return &Handler{
		client:     client,
		guest:      guestDescriptor,
		logger:     logger,
		dashboards: dashboards,
	}
}

func (h *Handler) RegisterSuperset(r *gin.RouterGroup) {
	sg := r.Group("/superset")
	sg.GET("/guest-token", h.GuestToken)
	sg.GET("/dashboards", h.Dashboards)
}

func (h *Handler) GuestToken(c *gin.Context) {
	dashboardId := c.Query("dashboard")

	token, err := h.client.GetGuestToken(dashboardId, h.guest)
	if err != nil {
		h.logger.Warn("failed to authenticate", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Dashboards(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"dashboards": h.dashboards})
}
