package repository

import (
	"github.com/Ruletk/OnlineClinic/apps/doctor/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DoctorRepository описывает доступ к хранилищу врачей
type DoctorRepository interface {
	Create(doc *model.Doctor) error
	GetByID(id uuid.UUID) (*model.Doctor, error)
	Update(doc *model.Doctor) error
	Delete(id uuid.UUID) error
}

type doctorRepo struct {
	db *gorm.DB
}

// NewDoctorRepository конструктор
func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepo{db: db}
}

func (r *doctorRepo) Create(doc *model.Doctor) error {
	return r.db.Create(doc).Error
}

func (r *doctorRepo) GetByID(id uuid.UUID) (*model.Doctor, error) {
	var doc model.Doctor
	if err := r.db.
		Preload("Specialization").
		Preload("ScheduleSlots").
		First(&doc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *doctorRepo) Update(doc *model.Doctor) error {
	return r.db.Save(doc).Error
}

func (r *doctorRepo) Delete(id uuid.UUID) error {
	return r.db.
		Delete(&model.Doctor{}, "id = ?", id).Error
}
