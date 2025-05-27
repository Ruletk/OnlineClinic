package handler

import (
	"github.com/Ruletk/OnlineClinic/apps/doctor/internal/model"
	"github.com/Ruletk/OnlineClinic/apps/doctor/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ScheduleSlotHandler struct {
	svc service.ScheduleSlotService
}

func NewScheduleSlotHandler(s service.ScheduleSlotService) *ScheduleSlotHandler {
	return &ScheduleSlotHandler{svc: s}
}

func (h *ScheduleSlotHandler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/slots")
	g.POST("", h.create)
	g.GET("/:id", h.getByID)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *ScheduleSlotHandler) create(c *gin.Context) {
	var req model.ScheduleSlot
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	out, err := h.svc.CreateSlot(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *ScheduleSlotHandler) getByID(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	out, err := h.svc.GetSlotByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *ScheduleSlotHandler) update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	var req model.ScheduleSlot
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = id
	out, err := h.svc.UpdateSlot(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *ScheduleSlotHandler) delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	if err := h.svc.DeleteSlot(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
