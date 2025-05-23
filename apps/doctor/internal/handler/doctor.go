package handler

import (
	"doctor/internal/model"
	"doctor/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	svc service.DoctorService
}

func NewDoctorHandler(s service.DoctorService) *DoctorHandler {
	return &DoctorHandler{svc: s}
}

// RegisterRoutes вешает HTTP-роуты на Gin-движок
func (h *DoctorHandler) RegisterRoutes(r *gin.Engine) {
	docs := r.Group("/doctors")
	{
		docs.POST("", h.Create)
		docs.GET("/:id", h.GetByID)
		docs.PUT("/:id", h.Update)
		docs.DELETE("/:id", h.Delete)
	}
}

func (h *DoctorHandler) Create(c *gin.Context) {
	var d model.Doctor
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Create(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func (h *DoctorHandler) GetByID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	d, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if d == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *DoctorHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var d model.Doctor
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d.ID = uint(id)
	if err := h.svc.Update(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (h *DoctorHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
