package dto

import (
	"github.com/google/uuid"
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
