package handler

import (
	"doctor/internal/model"
	"doctor/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SpecializationHandler struct {
	svc service.SpecializationService
}

func NewSpecializationHandler(s service.SpecializationService) *SpecializationHandler {
	return &SpecializationHandler{svc: s}
}

func (h *SpecializationHandler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/specializations")
	g.POST("", h.create)
	g.GET("/:id", h.getByID)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *SpecializationHandler) create(c *gin.Context) {
	var req model.Specialization
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.svc.CreateSpecialization(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *SpecializationHandler) getByID(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	out, err := h.svc.GetSpecializationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *SpecializationHandler) update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var req model.Specialization
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	out, err := h.svc.UpdateSpecialization(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *SpecializationHandler) delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.svc.DeleteSpecialization(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
