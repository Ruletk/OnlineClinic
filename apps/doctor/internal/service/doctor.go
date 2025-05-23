package service

import (
	"doctor/internal/model"
	"doctor/internal/repository"
	"errors"
)

type DoctorService interface {
	Create(d *model.Doctor) error
	GetByID(id uint) (*model.Doctor, error)
	Update(d *model.Doctor) error
	Delete(id uint) error
}

type doctorService struct {
	repo repository.DoctorRepository
}

func NewDoctorService(r repository.DoctorRepository) DoctorService {
	return &doctorService{repo: r}
}

func (s *doctorService) Create(d *model.Doctor) error {
	if d.Name == "" || d.Email == "" {
		return errors.New("name and email are required")
	}
	return s.repo.Create(d)
}

func (s *doctorService) GetByID(id uint) (*model.Doctor, error) {
	return s.repo.GetByID(id)
}

func (s *doctorService) Update(d *model.Doctor) error {
	if d.ID == 0 {
		return errors.New("id is required")
	}
	if d.Name == "" || d.Email == "" {
		return errors.New("name and email are required")
	}
	return s.repo.Update(d)
}

func (s *doctorService) Delete(id uint) error {
	if id == 0 {
		return errors.New("id is required")
	}
	return s.repo.Delete(id)
}
