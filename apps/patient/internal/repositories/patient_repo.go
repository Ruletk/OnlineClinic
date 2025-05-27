package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"patient/internal/models"
)

type PatientRepository interface {
	Create(patient *models.Patient) error
	GetByID(uuid.UUID) (*models.Patient, error)
	Update(patient *models.Patient) error
	Delete(uuid.UUID) error
	GetAll(i *[]models.Patient, limit int, offset int) error
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

func (r *PatientRepo) GetByID(id uuid.UUID) (*models.Patient, error) {
	var patient models.Patient
	if err := r.db.First(&patient, id).Error; err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *PatientRepo) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

func (r *PatientRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Patient{}, id).Error
}

func (r *PatientRepo) GetAll(patients *[]models.Patient, limit int, offset int) error {
	if limit <= 0 || offset < 0 {
		return r.db.Find(patients).Error
	}
	return r.db.Limit(limit).Offset(offset).Find(patients).Error
}
