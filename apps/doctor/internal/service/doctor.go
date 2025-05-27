package service

import (
	"context"
	"errors"
	"github.com/Ruletk/OnlineClinic/apps/doctor/internal/model"
	"github.com/Ruletk/OnlineClinic/apps/doctor/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrNotFound возвращается, если врач не найден
var ErrNotFound = errors.New("doctor not found")

// DoctorService описывает бизнес-логику над врачами
type DoctorService interface {
	CreateDoctor(ctx context.Context, req CreateDoctorRequest) (*CreateDoctorResponse, error)
	GetDoctorByID(ctx context.Context, id uuid.UUID) (*DoctorDTO, error)
	UpdateDoctor(ctx context.Context, req UpdateDoctorRequest) (*DoctorDTO, error)
	DeleteDoctor(ctx context.Context, id uuid.UUID) (*DeleteDoctorResponse, error)
}

type doctorService struct {
	repo repository.DoctorRepository
}

// NewDoctorService конструктор
func NewDoctorService(repo repository.DoctorRepository) DoctorService {
	return &doctorService{repo: repo}
}

func (s *doctorService) CreateDoctor(ctx context.Context, req CreateDoctorRequest) (*CreateDoctorResponse, error) {
	// Распакуем Patronymic (он у вас *string в DTO)
	pat := ""
	if req.Patronymic != nil {
		pat = *req.Patronymic
	}

	doc := &model.Doctor{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Patronymic:       pat,
		DateOfBirth:      req.DateOfBirth,
		SpecializationID: req.SpecializationID,
		// В ваших моделях константа называется Active, а не StatusActive
		Status: model.Active,
	}

	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, err
	}

	// Подготовим строковый указатель
	patPtr := doc.Patronymic

	return &CreateDoctorResponse{
		ID:               doc.ID,
		FirstName:        doc.FirstName,
		LastName:         doc.LastName,
		Patronymic:       &patPtr,
		DateOfBirth:      doc.DateOfBirth,
		SpecializationID: doc.SpecializationID,
		Status:           string(doc.Status),
		CreatedAt:        doc.CreatedAt,
	}, nil
}

func (s *doctorService) GetDoctorByID(ctx context.Context, id uuid.UUID) (*DoctorDTO, error) {
	doc, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Передадим указатель на Patronymic
	patPtr := doc.Patronymic

	return &DoctorDTO{
		ID:               doc.ID,
		FirstName:        doc.FirstName,
		LastName:         doc.LastName,
		Patronymic:       &patPtr,
		DateOfBirth:      doc.DateOfBirth,
		SpecializationID: doc.SpecializationID,
		Status:           string(doc.Status),
	}, nil
}

func (s *doctorService) UpdateDoctor(ctx context.Context, req UpdateDoctorRequest) (*DoctorDTO, error) {
	existing, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Обновляем поля
	existing.FirstName = req.FirstName
	existing.LastName = req.LastName
	if req.Patronymic != nil {
		existing.Patronymic = *req.Patronymic
	}
	existing.DateOfBirth = req.DateOfBirth
	existing.SpecializationID = req.SpecializationID
	existing.Status = model.DoctorStatus(req.Status)

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	patPtr := existing.Patronymic

	return &DoctorDTO{
		ID:               existing.ID,
		FirstName:        existing.FirstName,
		LastName:         existing.LastName,
		Patronymic:       &patPtr,
		DateOfBirth:      existing.DateOfBirth,
		SpecializationID: existing.SpecializationID,
		Status:           string(existing.Status),
	}, nil
}

func (s *doctorService) DeleteDoctor(ctx context.Context, id uuid.UUID) (*DeleteDoctorResponse, error) {
	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return &DeleteDoctorResponse{Success: true}, nil
}
