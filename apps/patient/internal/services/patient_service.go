package services

import (
	"patient/internal/dto"
	"patient/internal/repositories"
)

type PatientService interface {
	GetPatientByID(req *dto.GetPatientRequest) (*dto.PatientResponse, error)
	GetPatientAllergies(req *dto.GetPatientRequest) (*dto.AllergyResponses, error)
	GetPatientInsurances(req *dto.GetPatientRequest) (*dto.InsuranceResponses, error)
	GetPatientPrescriptions(req *dto.GetPatientRequest) (*dto.PrescriptionResponses, error)
	CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error)
	UpdatePatient(req *dto.UpdatePatientRequest) (*dto.PatientResponse, error)
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

func (p patientService) GetPatientByID(req *dto.GetPatientRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) GetPatientAllergies(req *dto.GetPatientRequest) (*dto.AllergyResponses, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) GetPatientInsurances(req *dto.GetPatientRequest) (*dto.InsuranceResponses, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) GetPatientPrescriptions(req *dto.GetPatientRequest) (*dto.PrescriptionResponses, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) UpdatePatient(req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) AddPatientAllergy(req *dto.CreateAllergyRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) AddPatientInsurance(req *dto.CreateInsuranceRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) AddPatientPrescription(req *dto.CreatePrescriptionRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) DeletePatientAllergy(req *dto.DeleteAllergyRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) DeletePatientInsurance(req *dto.DeleteInsuranceRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) DeletePatientPrescription(req *dto.DeletePrescriptionRequest) (*dto.PatientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p patientService) GetAllPatients(limit, offset int) (*dto.PatientResponses, error) {
	//TODO implement me
	panic("implement me")
}

func NewPatientService(
	repo repositories.PatientRepository,
	allergyService AllergyService,
	insuranceService InsuranceService,
	prescriptionService PrescriptionService,
) PatientService {
	return &patientService{
		repo:                repo,
		allergyService:      allergyService,
		insuranceService:    insuranceService,
		prescriptionService: prescriptionService,
	}
}
