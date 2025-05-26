package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreatePatientRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	BloodType string    `json:"blood_type" validate:"oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	Height    float64   `json:"height" validate:"min=0"`
	Weight    float64   `json:"weight" validate:"min=0"`
}

type GetPatientRequest struct {
	PatientID uuid.UUID `json:"patient_id" validate:"required"`
}

type UpdatePatientRequest struct {
	BloodType string  `json:"blood_type" validate:"omitempty,oneof=A+ A- B+ B- AB+ AB- O+ O-"`
	Height    float64 `json:"height" validate:"omitempty,min=0"`
	Weight    float64 `json:"weight" validate:"omitempty,min=0"`
}

type CreateAllergyRequest struct {
	PatientID  uuid.UUID `json:"patient_id" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	ObservedAt string    `json:"observed_at" validate:"required,datetime=2006-01-02"`
	Severity   string    `json:"severity" validate:"required,oneof=LOW MODERATE SEVERE"`
}

type DeleteAllergyRequest struct {
	AllergyID uuid.UUID `json:"allergy_id" validate:"required"`
}

type UpdateAllergyRequest struct {
	AllergyID  uuid.UUID `json:"allergy_id" validate:"required"`
	Name       string    `json:"name" validate:"omitempty"`
	ObservedAt string    `json:"observed_at" validate:"omitempty,datetime=2006-01-02"`
	Severity   string    `json:"severity" validate:"omitempty,oneof=LOW MODERATE SEVERE"`
}

type CreateInsuranceRequest struct {
	PatientID      uuid.UUID `json:"patient_id" validate:"required"`
	Provider       string    `json:"provider" validate:"required"`
	PolicyNumber   string    `json:"policy_number" validate:"required"`
	ExpirationDate string    `json:"expiration_date" validate:"required,datetime=2006-01-02"`
}

type UpdateInsuranceRequest struct {
	InsuranceID uuid.UUID `json:"insurance_id" validate:"required"`
}

type DeleteInsuranceRequest struct {
	InsuranceID uuid.UUID `json:"insurance_id" validate:"required"`
}

type CreatePrescriptionRequest struct {
	PatientID  uuid.UUID `json:"patient_id" validate:"required"`
	DoctorID   uuid.UUID `json:"doctor_id" validate:"required"`
	Medication string    `json:"medication" validate:"required"`
	Dosage     string    `json:"dosage" validate:"required"`
	ValidUntil time.Time `json:"valid_until" validate:"required,datetime=2006-01-02"`
}

type DeletePrescriptionRequest struct {
	PrescriptionID uuid.UUID `json:"prescription_id" validate:"required"`
}
