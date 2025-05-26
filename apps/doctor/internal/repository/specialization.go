package repository

import (
	"doctor/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpecializationRepository interface {
	Create(s *model.Specialization) error
	GetByID(id uuid.UUID) (*model.Specialization, error)
	Update(s *model.Specialization) error
	Delete(id uuid.UUID) error
}

type specializationRepo struct {
	db *gorm.DB
}

func NewSpecializationRepository(db *gorm.DB) SpecializationRepository {
	return &specializationRepo{db: db}
}

func (r *specializationRepo) Create(s *model.Specialization) error {
	return r.db.Create(s).Error
}

func (r *specializationRepo) GetByID(id uuid.UUID) (*model.Specialization, error) {
	var s model.Specialization
	if err := r.db.First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *specializationRepo) Update(s *model.Specialization) error {
	return r.db.Save(s).Error
}

func (r *specializationRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Specialization{}, "id = ?", id).Error
}
