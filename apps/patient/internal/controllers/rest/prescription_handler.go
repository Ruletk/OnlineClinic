package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"patient/internal/dto"
)

func (h *PatientHandler) RegisterPrescriptionRoutes(router *gin.Engine) {
	prescriptions := router.Group("/patients/:id/prescriptions")
	{
		prescriptions.POST("", h.AddPrescription)
		prescriptions.GET("", h.GetPrescriptions)
	}
	router.DELETE("/prescriptions/:id", h.DeletePrescription)
}

// AddPrescription - Добавить рецепт
func (h *PatientHandler) AddPrescription(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	var req dto.CreatePrescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.PatientID = patientID
	patient, err := h.service.AddPatientPrescription(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetPrescriptions - Получить рецепты пациента
func (h *PatientHandler) GetPrescriptions(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	prescriptions, err := h.service.GetPatientPrescriptions(&dto.GetPatientRequest{PatientID: patientID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prescriptions)
}

// DeletePrescription - Удалить рецепт
func (h *PatientHandler) DeletePrescription(c *gin.Context) {
	prescriptionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid prescription ID"})
		return
	}

	if err, _ := h.service.DeletePatientPrescription(&dto.DeletePrescriptionRequest{PrescriptionID: prescriptionID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
