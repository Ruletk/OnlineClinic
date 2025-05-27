package services

import (
	"errors"
	"github.com/google/uuid"
	"patient/internal/dto"
	"patient/internal/models"
	"patient/internal/repositories"
	"time"
)

type InsuranceService interface {
	Create(req *dto.CreateInsuranceRequest) (*dto.InsuranceResponse, error)
	Delete(req *dto.DeleteInsuranceRequest) error
}

type insuranceService struct {
	repo repositories.InsuranceRepository
}

func (i insuranceService) Create(req *dto.CreateInsuranceRequest) (*dto.InsuranceResponse, error) {
	// Parse ExpirationDate
	expirationDate, err := time.Parse("2006-01-02", req.ExpirationDate)
	if err != nil {
		return nil, errors.New("invalid date format for ExpirationDate")
	}

	// Map the request DTO to the model
	insurance := &models.Insurance{
		ID:             uuid.New(),
		PatientID:      req.PatientID,
		Provider:       req.Provider,
		PolicyNumber:   req.PolicyNumber,
		ExpirationDate: expirationDate,
	}

	// Save the insurance using the repository
	if err := i.repo.Create(insurance); err != nil {
		return nil, err
	}
	return &dto.InsuranceResponse{
		ID:             insurance.ID,
		PatientID:      insurance.PatientID,
		Provider:       insurance.Provider,
		PolicyNumber:   insurance.PolicyNumber,
		ExpirationDate: insurance.ExpirationDate,
	}, nil
}

func (i insuranceService) Delete(req *dto.DeleteInsuranceRequest) error {
	// Validate the ID
	if req.InsuranceID == uuid.Nil {
		return errors.New("invalid insurance ID")
	}

	// Call the repository to delete the insurance
	return i.repo.Delete(req.InsuranceID)
}

func NewInsuranceService(repo repositories.InsuranceRepository) InsuranceService {
	return &insuranceService{
		repo: repo,
	}
}
