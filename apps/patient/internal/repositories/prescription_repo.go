package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"patient/internal/models"
)

type PrescriptionRepository interface {
	Create(prescription *models.Prescription) error
	GetByID(uuid.UUID) (*models.Prescription, error)
	Update(prescription *models.Prescription) error
	Delete(uuid.UUID) error
	GetByPatient(uuid.UUID) ([]models.Prescription, error)
}

type PrescriptionRepo struct {
	db *gorm.DB
}

func (r *PrescriptionRepo) GetByPatient(u uuid.UUID) ([]models.Prescription, error) {
	var prescriptions []models.Prescription
	if err := r.db.Where("patient_id = ?", u).Find(&prescriptions).Error; err != nil {
		return nil, err
	}
	if len(prescriptions) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return prescriptions, nil
}

func NewPrescriptionRepository(db *gorm.DB) PrescriptionRepository {
	return &PrescriptionRepo{db: db}
}

func (r *PrescriptionRepo) Create(prescription *models.Prescription) error {
	return r.db.Create(prescription).Error
}

func (r *PrescriptionRepo) GetByID(id uuid.UUID) (*models.Prescription, error) {
	var prescription models.Prescription
	if err := r.db.First(&prescription, id).Error; err != nil {
		return nil, err
	}
	return &prescription, nil
}

func (r *PrescriptionRepo) Update(prescription *models.Prescription) error {
	return r.db.Save(prescription).Error
}

func (r *PrescriptionRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Prescription{}, id).Error
}
