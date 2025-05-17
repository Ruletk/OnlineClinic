package service

import (
	"OnlineClinic/apps/doctor/model"
	"OnlineClinic/apps/doctor/repository"
	"errors"
)

type DoctorService struct {
	repo repository.DoctorRepository
}

func NewDoctorService(repo repository.DoctorRepository) *DoctorService {
	return &DoctorService{repo: repo}
}

func (s *DoctorService) CreateDoctor(doctor *model.Doctor) error {
	if doctor.Name == "" || doctor.Email == "" {
		return errors.New("name and email are required")
	}
	return s.repo.Create(doctor)
}

func (s *DoctorService) GetDoctorByID(id int64) (*model.Doctor, error) {
	return s.repo.GetByID(id)
}

func (s *DoctorService) UpdateDoctor(doctor *model.Doctor) error {
	if doctor.ID == 0 {
		return errors.New("invalid doctor ID")
	}
	return s.repo.Update(doctor)
}

func (s *DoctorService) DeleteDoctor(id int64) error {
	return s.repo.Delete(id)
}
