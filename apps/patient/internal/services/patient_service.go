package services

import (
	"github.com/google/uuid"
	"patient/internal/dto"
	"patient/internal/models"
	"patient/internal/repositories"
)

type PatientService interface {
	GetPatientByID(req *dto.GetPatientRequest) (*dto.PatientResponse, error)
	GetPatientAllergies(req *dto.GetPatientRequest) (*dto.AllergyResponses, error)
	GetPatientInsurances(req *dto.GetPatientRequest) (*dto.InsuranceResponses, error)
	GetPatientPrescriptions(req *dto.GetPatientRequest) (*dto.PrescriptionResponses, error)
	CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error)
	UpdatePatient(req *dto.UpdatePatientRequest) (*dto.PatientResponse, error)
	DeletePatient(req uuid.UUID) *dto.PatientResponse
	AddPatientAllergy(req *dto.CreateAllergyRequest) (*dto.PatientResponse, error)
	AddPatientInsurance(req *dto.CreateInsuranceRequest) (*dto.PatientResponse, error)
	AddPatientPrescription(req *dto.CreatePrescriptionRequest) (*dto.PatientResponse, error)
	DeletePatientAllergy(req *dto.DeleteAllergyRequest) (*dto.PatientResponse, error)
	DeletePatientInsurance(req *dto.DeleteInsuranceRequest) (*dto.PatientResponse, error)
	DeletePatientPrescription(req *dto.DeletePrescriptionRequest) (*dto.PatientResponse, error)
	GetAllPatients(limit, offset int) (*dto.PatientResponses, error)
}

type patientService struct {
	repo repositories.PatientRepository

	allergyService      AllergyService
	insuranceService    InsuranceService
	prescriptionService PrescriptionService
}

func (p patientService) DeletePatient(req uuid.UUID) *dto.PatientResponse {
	//TODO implement me
	panic("implement me")
}

func (p patientService) GetPatientByID(req *dto.GetPatientRequest) (*dto.PatientResponse, error) {
	// Retrieve the patient by ID
	patient, err := p.repo.GetByID(req.PatientID)
	if err != nil {
		return nil, err
	}

	// Convert the patient to DTO
	return &dto.PatientResponse{
		ID:         patient.ID,
		UserID:     patient.UserID,
		BloodType:  patient.BloodType,
		Height:     patient.Height,
		Weight:     patient.Weight,
		Allergies:  *dto.NewAllergyResponses(patient.Allergies),
		Insurances: *dto.NewInsuranceResponses(patient.Insurances),
		CreatedAt:  patient.CreatedAt,
		UpdatedAt:  patient.UpdatedAt,
	}, nil
}

func (p patientService) GetPatientAllergies(req *dto.GetPatientRequest) (*dto.AllergyResponses, error) {
	// Retrieve the patient by ID
	patient, err := p.repo.GetByID(req.PatientID)
	if err != nil {
		return nil, err
	}

	// Convert allergies to DTO
	return dto.NewAllergyResponses(patient.Allergies), nil
}

func (p patientService) GetPatientInsurances(req *dto.GetPatientRequest) (*dto.InsuranceResponses, error) {
	// Retrieve the patient by ID
	patient, err := p.repo.GetByID(req.PatientID)
	if err != nil {
		return nil, err
	}

	// Convert insurances to DTO
	return dto.NewInsuranceResponses(patient.Insurances), nil
}

func (p patientService) GetPatientPrescriptions(req *dto.GetPatientRequest) (*dto.PrescriptionResponses, error) {
	// Retrieve prescriptions by patient ID
	prescriptions, err := p.prescriptionService.GetByPatient(req.PatientID)
	if err != nil {
		return nil, err
	}

	return prescriptions, nil
}

func (p patientService) CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error) {
	// Map the request DTO to the model
	patient := &models.Patient{
		UserID:    req.UserID,
		BloodType: req.BloodType,
		Height:    req.Height,
		Weight:    req.Weight,
	}

	// Save the patient using the repository
	if err := p.repo.Create(patient); err != nil {
		return nil, err
	}

	// Convert the patient to DTO
	return &dto.PatientResponse{
		ID:        patient.ID,
		UserID:    patient.UserID,
		BloodType: patient.BloodType,
		Height:    patient.Height,
		Weight:    patient.Weight,
		CreatedAt: patient.CreatedAt,
		UpdatedAt: patient.UpdatedAt,
	}, nil
}

func (p patientService) UpdatePatient(req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	// Retrieve the existing patient
	patient, err := p.repo.GetByID(req.PatientID)
	if err != nil {
		return nil, err
	}

	// Update the fields
	if req.BloodType != "" {
		patient.BloodType = req.BloodType
	}
	if req.Height > 0 {
		patient.Height = req.Height
	}
	if req.Weight > 0 {
		patient.Weight = req.Weight
	}

	// Save the updated patient
	if err := p.repo.Update(patient); err != nil {
		return nil, err
	}

	// Convert the updated patient to DTO
	return &dto.PatientResponse{
		ID:        patient.ID,
		UserID:    patient.UserID,
		BloodType: patient.BloodType,
		Height:    patient.Height,
		Weight:    patient.Weight,
		CreatedAt: patient.CreatedAt,
		UpdatedAt: patient.UpdatedAt,
	}, nil
}

func (p patientService) AddPatientAllergy(req *dto.CreateAllergyRequest) (*dto.PatientResponse, error) {
	// Создать аллергию через AllergyService
	_, err := p.allergyService.Create(req)
	if err != nil {
		return nil, err
	}

	// Получить обновлённого пациента
	patient, err := p.repo.GetByID(req.PatientID)
	if err != nil {
		return nil, err
	}

	// Преобразовать пациента в DTO
	return &dto.PatientResponse{
		ID:         patient.ID,
		UserID:     patient.UserID,
		BloodType:  patient.BloodType,
		Height:     patient.Height,
		Weight:     patient.Weight,
		Allergies:  *dto.NewAllergyResponses(patient.Allergies),
		Insurances: *dto.NewInsuranceResponses(patient.Insurances),
		CreatedAt:  patient.CreatedAt,
		UpdatedAt:  patient.UpdatedAt,
	}, nil
}

func (p patientService) AddPatientInsurance(req *dto.CreateInsuranceRequest) (*dto.PatientResponse, error) {
	// Создать страховку через InsuranceService
	_, err := p.insuranceService.Create(req)
	if err != nil {
		return nil, err
	}

	// Получить обновлённого пациента
	patient, err := p.repo.GetByID(req.PatientID)
	if err != nil {
		return nil, err
	}

	// Преобразовать пациента в DTO
	return &dto.PatientResponse{
		ID:         patient.ID,
		UserID:     patient.UserID,
		BloodType:  patient.BloodType,
		Height:     patient.Height,
		Weight:     patient.Weight,
		Allergies:  *dto.NewAllergyResponses(patient.Allergies),
		Insurances: *dto.NewInsuranceResponses(patient.Insurances),
		CreatedAt:  patient.CreatedAt,
		UpdatedAt:  patient.UpdatedAt,
	}, nil
}

func (p patientService) AddPatientPrescription(req *dto.CreatePrescriptionRequest) (*dto.PatientResponse, error) {
	// Создать прескрипшон через PrescriptionService
	_, err := p.prescriptionService.Create(req)
	if err != nil {
		return nil, err
	}

	// Получить обновлённого пациента
	patient, err := p.repo.GetByID(req.PatientID)
	if err != nil {
		return nil, err
	}

	// Преобразовать пациента в DTO
	return &dto.PatientResponse{
		ID:         patient.ID,
		UserID:     patient.UserID,
		BloodType:  patient.BloodType,
		Height:     patient.Height,
		Weight:     patient.Weight,
		Allergies:  *dto.NewAllergyResponses(patient.Allergies),
		Insurances: *dto.NewInsuranceResponses(patient.Insurances),
		CreatedAt:  patient.CreatedAt,
		UpdatedAt:  patient.UpdatedAt,
	}, nil
}

func (p patientService) DeletePatientAllergy(req *dto.DeleteAllergyRequest) (*dto.PatientResponse, error) {
	// Delete the allergy through AllergyService
	err := p.allergyService.Delete(req)
	if err != nil {
		return nil, err
	}

	// Retrieve the updated patient
	patient, err := p.repo.GetByID(req.AllergyID)
	if err != nil {
		return nil, err
	}

	// Convert the patient to DTO
	return &dto.PatientResponse{
		ID:         patient.ID,
		UserID:     patient.UserID,
		BloodType:  patient.BloodType,
		Height:     patient.Height,
		Weight:     patient.Weight,
		Allergies:  *dto.NewAllergyResponses(patient.Allergies),
		Insurances: *dto.NewInsuranceResponses(patient.Insurances),
		CreatedAt:  patient.CreatedAt,
		UpdatedAt:  patient.UpdatedAt,
	}, nil
}

func (p patientService) DeletePatientInsurance(req *dto.DeleteInsuranceRequest) (*dto.PatientResponse, error) {
	// Delete the insurance through InsuranceService
	err := p.insuranceService.Delete(req)
	if err != nil {
		return nil, err
	}

	// Retrieve the updated patient
	patient, err := p.repo.GetByID(req.InsuranceID)
	if err != nil {
		return nil, err
	}

	// Convert the patient to DTO
	return &dto.PatientResponse{
		ID:         patient.ID,
		UserID:     patient.UserID,
		BloodType:  patient.BloodType,
		Height:     patient.Height,
		Weight:     patient.Weight,
		Allergies:  *dto.NewAllergyResponses(patient.Allergies),
		Insurances: *dto.NewInsuranceResponses(patient.Insurances),
		CreatedAt:  patient.CreatedAt,
		UpdatedAt:  patient.UpdatedAt,
	}, nil
}

func (p patientService) DeletePatientPrescription(req *dto.DeletePrescriptionRequest) (*dto.PatientResponse, error) {
	// Delete the prescription through PrescriptionService
	err := p.prescriptionService.Delete(req)
	if err != nil {
		return nil, err
	}

	// Retrieve the updated patient
	patient, err := p.repo.GetByID(req.PrescriptionID)
	if err != nil {
		return nil, err
	}

	// Convert the patient to DTO
	return &dto.PatientResponse{
		ID:         patient.ID,
		UserID:     patient.UserID,
		BloodType:  patient.BloodType,
		Height:     patient.Height,
		Weight:     patient.Weight,
		Allergies:  *dto.NewAllergyResponses(patient.Allergies),
		Insurances: *dto.NewInsuranceResponses(patient.Insurances),
		CreatedAt:  patient.CreatedAt,
		UpdatedAt:  patient.UpdatedAt,
	}, nil
}

func (p patientService) GetAllPatients(limit, offset int) (*dto.PatientResponses, error) {
	// Retrieve all patients with pagination
	var patients []models.Patient
	err := p.repo.GetAll(&patients, limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert patients to DTO
	responses := make([]dto.PatientResponse, len(patients))
	for i, patient := range patients {
		responses[i] = dto.PatientResponse{
			ID:         patient.ID,
			UserID:     patient.UserID,
			BloodType:  patient.BloodType,
			Height:     patient.Height,
			Weight:     patient.Weight,
			Allergies:  *dto.NewAllergyResponses(patient.Allergies),
			Insurances: *dto.NewInsuranceResponses(patient.Insurances),
			CreatedAt:  patient.CreatedAt,
			UpdatedAt:  patient.UpdatedAt,
		}
	}

	return &dto.PatientResponses{
		Count:    len(responses),
		Patients: responses,
	}, nil
}

func NewPatientService(repo repositories.PatientRepository, allergyService *allergyService, insuranceService InsuranceService, prescriptionService PrescriptionService) PatientService {
	return &patientService{
		repo:                repo,
		allergyService:      allergyService,
		insuranceService:    insuranceService,
		prescriptionService: prescriptionService,
	}
}
