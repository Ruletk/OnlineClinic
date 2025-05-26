package repository

import (
	"context"
	"doctor/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DoctorRepository описывает доступ к хранилищу врачей
type DoctorRepository interface {
	Create(ctx context.Context, doc *model.Doctor) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Doctor, error)
	Update(ctx context.Context, doc *model.Doctor) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type doctorRepo struct {
	db *gorm.DB
}

// NewDoctorRepository конструктор
func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepo{db: db}
}

func (r *doctorRepo) Create(ctx context.Context, doc *model.Doctor) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *doctorRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Doctor, error) {
	var doc model.Doctor
	if err := r.db.WithContext(ctx).
		Preload("Specialization").
		Preload("ScheduleSlots").
		First(&doc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *doctorRepo) Update(ctx context.Context, doc *model.Doctor) error {
	return r.db.WithContext(ctx).Save(doc).Error
}

func (r *doctorRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&model.Doctor{}, "id = ?", id).Error
}
