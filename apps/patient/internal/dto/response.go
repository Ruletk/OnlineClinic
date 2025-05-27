package dto

import (
	"github.com/google/uuid"
	"patient/internal/models"
	"time"
)

type PatientResponse struct {
	ID         int64              `json:"id"`
	UserID     uuid.UUID          `json:"user_id"`
	BloodType  string             `json:"blood_type"`
	Height     float64            `json:"height"`
	Weight     float64            `json:"weight"`
	Allergies  AllergyResponses   `json:"allergies"`
	Insurances InsuranceResponses `json:"insurances"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

func (r PatientResponse) Error() any {
	return map[string]any{
		"error": "Patient not found",
	}
}

type PatientResponses struct {
	Count    int               `json:"count"`
	Patients []PatientResponse `json:"patients"`
}

type AllergyResponse struct {
	ID         uuid.UUID `json:"id"`
	PatientID  uuid.UUID `json:"patient_id"`
	Name       string    `json:"name"`
	Severity   string    `json:"severity"`
	ObservedAt time.Time `json:"observed_at"`
}

type AllergyResponses struct {
	Count     int               `json:"count"`
	Allergies []AllergyResponse `json:"allergies"`
}

type InsuranceResponse struct {
	ID             uuid.UUID `json:"id"`
	PatientID      uuid.UUID `json:"patient_id"`
	Provider       string    `json:"provider"`
	PolicyNumber   string    `json:"policy_number"`
	ExpirationDate time.Time `json:"expiration_date"`
}

type InsuranceResponses struct {
	Count      int                 `json:"count"`
	Insurances []InsuranceResponse `json:"insurances"`
}

type PrescriptionResponse struct {
	ID         uuid.UUID `json:"id"`
	PatientID  uuid.UUID `json:"patient_id"`
	DoctorID   uuid.UUID `json:"doctor_id"`
	Medication string    `json:"medication"`
	Dosage     string    `json:"dosage"`
	ValidUntil time.Time `json:"valid_until"`
}

type PrescriptionResponses struct {
	Count         int                    `json:"count"`
	Prescriptions []PrescriptionResponse `json:"prescriptions"`
}

func NewPrescriptionResponse(prescription *models.Prescription) *PrescriptionResponse {
	return &PrescriptionResponse{
		ID:         prescription.ID,
		PatientID:  prescription.PatientID,
		DoctorID:   prescription.DoctorID,
		Medication: prescription.Medication,
		Dosage:     prescription.Dosage,
		ValidUntil: prescription.ValidUntil,
	}
}

func NewPrescriptionResponses(prescriptions []models.Prescription) *PrescriptionResponses {
	responses := make([]PrescriptionResponse, len(prescriptions))
	for i, prescription := range prescriptions {
		responses[i] = *NewPrescriptionResponse(&prescription)
	}
	return &PrescriptionResponses{
		Count:         len(responses),
		Prescriptions: responses,
	}
}

func NewInsuranceResponses(insurances []models.Insurance) *InsuranceResponses {
	responses := make([]InsuranceResponse, len(insurances))
	for i, insurance := range insurances {
		responses[i] = InsuranceResponse{
			ID:             insurance.ID,
			PatientID:      insurance.PatientID,
			Provider:       insurance.Provider,
			PolicyNumber:   insurance.PolicyNumber,
			ExpirationDate: insurance.ExpirationDate,
		}
	}
	return &InsuranceResponses{
		Count:      len(responses),
		Insurances: responses,
	}

}

func NewAllergyResponses(allergies []models.Allergy) *AllergyResponses {
	responses := make([]AllergyResponse, len(allergies))
	for i, allergy := range allergies {
		responses[i] = AllergyResponse{
			ID:         allergy.ID,
			PatientID:  allergy.PatientID,
			Name:       allergy.Name,
			Severity:   allergy.Severity,
			ObservedAt: allergy.ObservedAt,
		}
	}
	return &AllergyResponses{
		Count:     len(responses),
		Allergies: responses,
	}

}
