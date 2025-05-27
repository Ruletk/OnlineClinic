package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"patient/internal/dto"
	"patient/internal/services"
	"strconv"
)

type PatientHandler struct {
	service services.PatientService
}

func NewPatientHandler(service services.PatientService) *PatientHandler {
	return &PatientHandler{service: service}
}

func (h *PatientHandler) RegisterRoutes(router *gin.Engine) {
	patients := router.Group("/patients")
	{
		patients.GET("", h.GetAllPatients)
		patients.POST("", h.CreatePatient)
		patients.GET("/:id", h.GetPatientByID)
		patients.PATCH("/:id", h.UpdatePatient)
		patients.DELETE("/:id", h.DeletePatient)
	}
}

// GetAllPatients - Получить всех пациентов
func (h *PatientHandler) GetAllPatients(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	patients, err := h.service.GetAllPatients(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patients)
}

// CreatePatient - Создать пациента
func (h *PatientHandler) CreatePatient(c *gin.Context) {
	var req dto.CreatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patient, err := h.service.CreatePatient(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetPatientByID - Получить пациента по ID
func (h *PatientHandler) GetPatientByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	patient, err := h.service.GetPatientByID(&dto.GetPatientRequest{PatientID: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// UpdatePatient - Обновить данные пациента
func (h *PatientHandler) UpdatePatient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	var req dto.UpdatePatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.PatientID = id
	patient, err := h.service.UpdatePatient(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

// DeletePatient - Удалить пациента
func (h *PatientHandler) DeletePatient(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	if err := h.service.DeletePatient(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
