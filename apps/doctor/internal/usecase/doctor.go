package usecase

import (
	"context"
	"doctor/internal/model"
	"doctor/internal/repository"
	"errors"
	"strings"
)

// UseCase defines business logic for doctor service
type UseCase interface {
	Create(ctx context.Context, doctor *model.Doctor) error
	GetByID(ctx context.Context, id int64) (*model.Doctor, error)
	Update(ctx context.Context, doctor *model.Doctor) error
	Delete(ctx context.Context, id int64) error
}

type useCase struct {
	repo repository.Repository
}

// NewUseCase initializes and returns the UseCase implementation
func NewUseCase(repo repository.Repository) UseCase {
	return &useCase{repo: repo}
}

// Create validates and creates a new doctor record
func (u *useCase) Create(ctx context.Context, doctor *model.Doctor) error {
	if err := validateDoctor(doctor); err != nil {
		return err
	}
	return u.repo.Create(ctx, doctor)
}

// GetByID retrieves a doctor by ID with basic validation
func (u *useCase) GetByID(ctx context.Context, id int64) (*model.Doctor, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID: must be > 0")
	}
	return u.repo.GetByID(ctx, id)
}

// Update validates and updates an existing doctor
func (u *useCase) Update(ctx context.Context, doctor *model.Doctor) error {
	if doctor.ID <= 0 {
		return errors.New("invalid ID for update")
	}
	if err := validateDoctor(doctor); err != nil {
		return err
	}
	return u.repo.Update(ctx, doctor)
}

// Delete removes a doctor by ID
func (u *useCase) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid ID for delete")
	}
	return u.repo.Delete(ctx, id)
}

// validateDoctor performs basic field checks
func validateDoctor(d *model.Doctor) error {
	if strings.TrimSpace(d.Name) == "" {
		return errors.New("doctor name cannot be empty")
	}
	if strings.TrimSpace(d.Specialty) == "" {
		return errors.New("doctor specialty cannot be empty")
	}
	if !strings.Contains(d.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}
