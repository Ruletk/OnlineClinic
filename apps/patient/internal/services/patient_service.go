package services

import (
	"patient/internal/models"
	"patient/internal/repositories"
)

type PatientService struct {
	repo repositories.PatientRepository
}

func NewPatientService(repo repositories.PatientRepository) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) CreatePatient(patient *models.Patient) error {
	return s.repo.Create(patient)
}

func (s *PatientService) GetPatient(id uint) (*models.Patient, error) {
	return s.repo.GetByID(id)
}

func (s *PatientService) UpdatePatient(patient *models.Patient) error {
	return s.repo.Update(patient)
}

func (s *PatientService) DeletePatient(id uint) error {
	return s.repo.Delete(id)
}
