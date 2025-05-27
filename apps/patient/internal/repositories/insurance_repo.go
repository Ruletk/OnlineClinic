package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"patient/internal/models"
)

type InsuranceRepository interface {
	Create(insurance *models.Insurance) error
	GetByID(uuid.UUID) (*models.Insurance, error)
	Update(insurance *models.Insurance) error
	Delete(uuid.UUID) error
}

type InsuranceRepo struct {
	db *gorm.DB
}

func NewInsuranceRepository(db *gorm.DB) InsuranceRepository {
	return &InsuranceRepo{db: db}
}

func (r *InsuranceRepo) Create(insurance *models.Insurance) error {
	return r.db.Create(insurance).Error
}

func (r *InsuranceRepo) GetByID(id uuid.UUID) (*models.Insurance, error) {
	var insurance models.Insurance
	if err := r.db.First(&insurance, id).Error; err != nil {
		return nil, err
	}
	return &insurance, nil
}

func (r *InsuranceRepo) Update(insurance *models.Insurance) error {
	return r.db.Save(insurance).Error

}

func (r *InsuranceRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Insurance{}, id).Error
}
