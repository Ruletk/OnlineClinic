package services

import (
	"github.com/google/uuid"
	"patient/internal/dto"
	"patient/internal/models"
	"patient/internal/repositories"
)

type PrescriptionService interface {
	Create(req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error)
	Delete(req *dto.DeletePrescriptionRequest) error
	GetByPatient(id uuid.UUID) (*dto.PrescriptionResponses, error)
}

type prescriptionService struct {
	repo repositories.PrescriptionRepository
}

func (p prescriptionService) Create(req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	prescription := &models.Prescription{
		ID:         uuid.New(),
		PatientID:  req.PatientID,
		DoctorID:   req.DoctorID,
		Medication: req.Medication,
		Dosage:     req.Dosage,
		ValidUntil: req.ValidUntil,
	}

	if err := p.repo.Create(prescription); err != nil {
		return nil, err
	}

	return &dto.PrescriptionResponse{
		ID:         prescription.ID,
		PatientID:  prescription.PatientID,
		DoctorID:   prescription.DoctorID,
		Medication: prescription.Medication,
		Dosage:     prescription.Dosage,
		ValidUntil: prescription.ValidUntil,
	}, nil
}

func (p prescriptionService) Delete(req *dto.DeletePrescriptionRequest) error {
	return p.repo.Delete(req.PrescriptionID)

}

func (p prescriptionService) GetByPatient(id uuid.UUID) (*dto.PrescriptionResponses, error) {
	prescriptions, err := p.repo.GetByPatient(id)

	if err != nil {
		return nil, err
	}
	responses := dto.NewPrescriptionResponses(prescriptions)
	return responses, nil
}

func NewPrescriptionService(repo repositories.PrescriptionRepository) PrescriptionService {
	return &prescriptionService{
		repo: repo,
	}

}
