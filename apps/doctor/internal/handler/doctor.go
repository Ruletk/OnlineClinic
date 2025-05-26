package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"doctor/internal/service"
)

type DoctorHandler struct {
	svc service.DoctorService
}

func NewDoctorHandler(s service.DoctorService) *DoctorHandler {
	return &DoctorHandler{svc: s}
}

// RegisterRoutes навешивает CRUD-роуты
func (h *DoctorHandler) RegisterRoutes(r *gin.Engine) {
	docs := r.Group("/doctors")
	{
		docs.POST("", h.CreateDoctor)
		docs.GET("/:id", h.GetDoctorByID)
		docs.PUT("/:id", h.UpdateDoctor)
		docs.DELETE("/:id", h.DeleteDoctor)
	}
}

// CreateDoctor — POST /doctors
func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var req service.CreateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.CreateDoctor(c.Request.Context(), req)
	if err != nil {
		// можно разделить типы ошибок и вернуть 409/422 и т.п.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetDoctorByID — GET /doctors/:id
func (h *DoctorHandler) GetDoctorByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	dto, err := h.svc.GetDoctorByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, dto)
}

// UpdateDoctor — PUT /doctors/:id
func (h *DoctorHandler) UpdateDoctor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req service.UpdateDoctorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id

	dto, err := h.svc.UpdateDoctor(c.Request.Context(), req)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, dto)
}

// DeleteDoctor — DELETE /doctors/:id
func (h *DoctorHandler) DeleteDoctor(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	_, err = h.svc.DeleteDoctor(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
