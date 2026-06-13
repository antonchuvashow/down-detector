package apiinspector

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"detector/internal/inspector/application"
	route "detector/internal/route/domain"
)

type Handler struct {
	bridge *inspectorapp.RouteInspectorBridge
}

func NewHandler(bridge *inspectorapp.RouteInspectorBridge) *Handler {
	return &Handler{
		bridge: bridge,
	}
}

func (h *Handler) RegisterInspectors(r *gin.RouterGroup) {
	inspectors := r.Group("/routes/:id/inspector")

	inspectors.GET("", h.Get)
	inspectors.PUT("", h.CreateOrUpdate)
	inspectors.DELETE("", h.Delete)
}

func (h *Handler) Get(c *gin.Context) {
	routeId := c.Param("id")
	inspectorInstance, err := h.bridge.FindInspector(route.ID(routeId))
	if err != nil {
		if _, ok := errors.AsType[*inspectorapp.ErrNotFound](err); ok {
			c.Status(http.StatusNotFound)
			return
		}
		c.Status(http.StatusInternalServerError)
		return
	}
	apiInspector, err := InspectorFromDomain(inspectorInstance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, apiInspector)
}

func (h *Handler) CreateOrUpdate(c *gin.Context) {
	routeId := route.ID(c.Param("id"))

	var createRequest CreateRequest
	if err := c.ShouldBindJSON(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	inspectorInstance, err := InspectorToDomain(createRequest.Inspector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = h.bridge.Register(routeId, inspectorInstance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	routeId := route.ID(c.Param("id"))
	err := h.bridge.DeleteInspector(routeId)
	if err != nil {
		if _, ok := errors.AsType[*inspectorapp.ErrNotFound](err); ok {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
