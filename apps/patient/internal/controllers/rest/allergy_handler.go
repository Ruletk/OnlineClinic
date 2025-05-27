package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"patient/internal/dto"
)

func (h *PatientHandler) RegisterAllergyRoutes(router *gin.Engine) {
	allergies := router.Group("/patients/:id/allergies")
	{
		allergies.POST("", h.AddAllergy)
		allergies.GET("", h.GetAllergies)
	}
	router.DELETE("/allergies/:id", h.DeleteAllergy)
}

// AddAllergy - Добавить аллергию
func (h *PatientHandler) AddAllergy(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	var req dto.CreateAllergyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.PatientID = patientID
	patient, err := h.service.AddPatientAllergy(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetAllergies - Получить аллергии пациента
func (h *PatientHandler) GetAllergies(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	allergies, err := h.service.GetPatientAllergies(&dto.GetPatientRequest{PatientID: patientID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, allergies)
}

// DeleteAllergy - Удалить аллергию
func (h *PatientHandler) DeleteAllergy(c *gin.Context) {
	allergyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid allergy ID"})
		return
	}

	if _, err := h.service.DeletePatientAllergy(&dto.DeleteAllergyRequest{AllergyID: allergyID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
