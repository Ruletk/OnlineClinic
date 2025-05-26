package repositories

import (
	"gorm.io/gorm"
	"patient/internal/models"
)

type AllergyRepository interface {
	Create(allergy *models.Allergy) error
	GetByID(id uint) (*models.Allergy, error)
	Update(allergy *models.Allergy) error
	Delete(id uint) error
}

type AllergyRepo struct {
	db *gorm.DB
}

func NewAllergyRepository(db *gorm.DB) AllergyRepository {
	return &AllergyRepo{db: db}
}

func (r *AllergyRepo) Create(allergy *models.Allergy) error {
	return r.db.Create(allergy).Error
}

func (r *AllergyRepo) GetByID(id uint) (*models.Allergy, error) {
	var allergy models.Allergy
	if err := r.db.First(&allergy, id).Error; err != nil {
		return nil, err
	}
	return &allergy, nil
}

func (r *AllergyRepo) Update(allergy *models.Allergy) error {
	return r.db.Save(allergy).Error
}

func (r *AllergyRepo) Delete(id uint) error {
	return r.db.Delete(&models.Allergy{}, id).Error
}
