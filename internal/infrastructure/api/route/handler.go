package apiroute

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"detector/internal/route/application"
	"detector/internal/route/domain"
)

type Handler struct {
	service *routeapp.Service
}

func NewHandler(service *routeapp.Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all route-related endpoints
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	routes := r.Group("/routes")
	{
		routes.GET("", h.List)
		routes.POST("", h.Create)
		routes.GET("/:id", h.Get)
		routes.PUT("/:id", h.Update)
		routes.DELETE("/:id", h.Delete)
	}
}

func (h *Handler) List(c *gin.Context) {
	routes, err := h.service.GetAllRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, NewListResponse(routes))
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse(err))
		return
	}

	result, err := h.service.Add(req.ToCommand())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, NewResponse(result))
}

func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.Get(route.ID(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := errors.AsType[routeapp.ErrNotFound](err); ok {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, NewResponse(result))
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, ValidationErrorResponse(err))
		return
	}

	err := h.service.Update(req.ToCommand(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := errors.AsType[routeapp.ErrNotFound](err); ok {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(route.ID(id)); err != nil {
		statusCode := http.StatusInternalServerError
		if _, ok := errors.AsType[routeapp.ErrNotFound](err); ok {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
