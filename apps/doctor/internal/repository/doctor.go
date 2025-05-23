package repository

import (
	"doctor/internal/model"
	"errors"

	"gorm.io/gorm"
)

// DoctorRepository описывает доступ к данным doctors
type DoctorRepository interface {
	Create(d *model.Doctor) error
	GetByID(id uint) (*model.Doctor, error)
	Update(d *model.Doctor) error
	Delete(id uint) error
}

type doctorRepo struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepo{db: db}
}

func (r *doctorRepo) Create(d *model.Doctor) error {
	return r.db.Create(d).Error
}

func (r *doctorRepo) GetByID(id uint) (*model.Doctor, error) {
	var d model.Doctor
	if err := r.db.First(&d, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *doctorRepo) Update(d *model.Doctor) error {
	return r.db.Save(d).Error
}

func (r *doctorRepo) Delete(id uint) error {
	return r.db.Delete(&model.Doctor{}, id).Error
}
