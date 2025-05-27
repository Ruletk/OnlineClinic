package services

import (
	"errors"
	"github.com/google/uuid"
	"patient/internal/dto"
	"patient/internal/models"
	"patient/internal/repositories"
	"time"
)

type AllergyService interface {
	Create(req *dto.CreateAllergyRequest) (*dto.AllergyResponse, error)
	Update(req *dto.UpdateAllergyRequest) (*dto.AllergyResponse, error)
	Delete(req *dto.DeleteAllergyRequest) error
}

type allergyService struct {
	repo repositories.AllergyRepository
}

func (a allergyService) Create(req *dto.CreateAllergyRequest) (*dto.AllergyResponse, error) {
	// Parse ObservedAt date
	observedAt, err := time.Parse("2006-01-02", req.ObservedAt)
	if err != nil {
		return nil, errors.New("invalid date format for ObservedAt")
	}

	// Map the request DTO to the model
	allergy := &models.Allergy{
		ID:         uuid.New(),
		PatientID:  req.PatientID,
		Name:       req.Name,
		Severity:   req.Severity,
		ObservedAt: observedAt,
	}

	// Save the allergy using the repository
	if err := a.repo.Create(allergy); err != nil {
		return nil, err
	}

	// Map the model back to the response DTO
	return &dto.AllergyResponse{
		ID:         allergy.ID,
		PatientID:  allergy.PatientID,
		Name:       allergy.Name,
		Severity:   allergy.Severity,
		ObservedAt: allergy.ObservedAt,
	}, nil
}

func (a allergyService) Update(req *dto.UpdateAllergyRequest) (*dto.AllergyResponse, error) {
	// Retrieve the existing allergy
	allergy, err := a.repo.GetByID(req.AllergyID)
	if err != nil {
		return nil, err
	}

	// Update the fields
	if req.Name != "" {
		allergy.Name = req.Name
	}
	if req.Severity != "" {
		allergy.Severity = req.Severity
	}
	if req.ObservedAt != "" {
		observedAt, err := time.Parse("2006-01-02", req.ObservedAt)
		if err != nil {
			return nil, errors.New("invalid date format for ObservedAt")
		}
		allergy.ObservedAt = observedAt
	}

	// Save the updated allergy
	if err := a.repo.Update(allergy); err != nil {
		return nil, err
	}

	// Map the updated model back to the response DTO
	return &dto.AllergyResponse{
		ID:         allergy.ID,
		PatientID:  allergy.PatientID,
		Name:       allergy.Name,
		Severity:   allergy.Severity,
		ObservedAt: allergy.ObservedAt,
	}, nil
}

func (a allergyService) Delete(req *dto.DeleteAllergyRequest) error {
	// Validate the ID
	if req.AllergyID == uuid.Nil {
		return errors.New("invalid allergy ID")
	}

	// Call the repository to delete the allergy
	return a.repo.Delete(req.AllergyID)
}

func NewAllergyService(repo repositories.AllergyRepository) AllergyService {
	return &allergyService{
		repo: repo,
	}
}
