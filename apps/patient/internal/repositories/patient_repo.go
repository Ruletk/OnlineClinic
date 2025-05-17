package repositories

import (
	"gorm.io/gorm"
	"patient/internal/models"
)

type PatientRepository interface {
	Create(patient *models.Patient) error
	GetByID(id uint) (*models.Patient, error)
	Update(patient *models.Patient) error
	Delete(id uint) error
}

type PatientRepo struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &PatientRepo{db: db}
}

func (r *PatientRepo) Create(patient *models.Patient) error {
	return r.db.Create(patient).Error
}

func (r *PatientRepo) GetByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.First(&patient, id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepo) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

func (r *PatientRepo) Delete(id uint) error {
	return r.db.Delete(&models.Patient{}, id).Error
}
