package service

import (
	"doctor/internal/model"
	"doctor/internal/repository"
	"errors"

	"github.com/google/uuid"
)

var ErrSpecNotFound = errors.New("specialization not found")

type SpecializationService interface {
	CreateSpecialization(req model.Specialization) (*model.Specialization, error)
	GetSpecializationByID(id uuid.UUID) (*model.Specialization, error)
	UpdateSpecialization(req model.Specialization) (*model.Specialization, error)
	DeleteSpecialization(id uuid.UUID) error
}

type specializationService struct {
	repo repository.SpecializationRepository
}

func NewSpecializationService(r repository.SpecializationRepository) SpecializationService {
	return &specializationService{repo: r}
}

func (s *specializationService) CreateSpecialization(req model.Specialization) (*model.Specialization, error) {
	req.ID = uuid.New()
	if err := s.repo.Create(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *specializationService) GetSpecializationByID(id uuid.UUID) (*model.Specialization, error) {
	spec, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrSpecNotFound
	}
	return spec, nil
}

func (s *specializationService) UpdateSpecialization(req model.Specialization) (*model.Specialization, error) {
	// ensure exists
	if _, err := s.repo.GetByID(req.ID); err != nil {
		return nil, ErrSpecNotFound
	}
	if err := s.repo.Update(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (s *specializationService) DeleteSpecialization(id uuid.UUID) error {
	return s.repo.Delete(id)
}
