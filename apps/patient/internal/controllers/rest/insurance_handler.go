package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"patient/internal/dto"
)

func (h *PatientHandler) RegisterInsuranceRoutes(router *gin.Engine) {
	insurances := router.Group("/patients/:id/insurances")
	{
		insurances.POST("", h.AddInsurance)
		insurances.GET("", h.GetInsurances)
	}
	router.DELETE("/insurances/:id", h.DeleteInsurance)
}

// AddInsurance - Добавить страховку
func (h *PatientHandler) AddInsurance(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	var req dto.CreateInsuranceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.PatientID = patientID
	patient, err := h.service.AddPatientInsurance(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, patient)
}

// GetInsurances - Получить страховки пациента
func (h *PatientHandler) GetInsurances(c *gin.Context) {
	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient ID"})
		return
	}

	insurances, err := h.service.GetPatientInsurances(&dto.GetPatientRequest{PatientID: patientID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insurances)
}

// DeleteInsurance - Удалить страховку
func (h *PatientHandler) DeleteInsurance(c *gin.Context) {
	insuranceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid insurance ID"})
		return
	}

	if err, _ := h.service.DeletePatientInsurance(&dto.DeleteInsuranceRequest{InsuranceID: insuranceID}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
