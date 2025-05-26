package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"patient/internal/models"
)

type AllergyRepository interface {
	Create(allergy *models.Allergy) error
	GetByID(uuid.UUID) (*models.Allergy, error)
	Update(allergy *models.Allergy) error
	Delete(uuid.UUID) error
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

func (r *AllergyRepo) GetByID(id uuid.UUID) (*models.Allergy, error) {
	var allergy models.Allergy
	if err := r.db.First(&allergy, id).Error; err != nil {
		return nil, err
	}
	return &allergy, nil
}

func (r *AllergyRepo) Update(allergy *models.Allergy) error {
	return r.db.Save(allergy).Error
}

func (r *AllergyRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Allergy{}, id).Error
}
