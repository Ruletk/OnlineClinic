package repositories

import (
	"gorm.io/gorm"
	"patient/internal/models"
)

type InsuranceRepository interface {
	Create(insurance *models.Insurance) error
	GetByID(id uint) (*models.Insurance, error)
	Update(insurance *models.Insurance) error
	Delete(id uint) error
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

func (r *InsuranceRepo) GetByID(id uint) (*models.Insurance, error) {
	var insurance models.Insurance
	if err := r.db.First(&insurance, id).Error; err != nil {
		return nil, err
	}
	return &insurance, nil
}

func (r *InsuranceRepo) Update(insurance *models.Insurance) error {
	return r.db.Save(insurance).Error

}

func (r *InsuranceRepo) Delete(id uint) error {
	return r.db.Delete(&models.Insurance{}, id).Error
}
