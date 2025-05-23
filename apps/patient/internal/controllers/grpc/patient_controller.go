package grpcpackage

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"patient/internal/models"
	proto "patient/internal/proto/gen"
	"patient/internal/services"
)

type PatientController struct {
	proto.UnimplementedPatientServiceServer
	service *services.PatientService
}

func NewPatientController(svc *services.PatientService) *PatientController {
	return &PatientController{service: svc}
}

func (c *PatientController) CreatePatient(ctx context.Context, req *proto.PatientRequest) (*proto.PatientResponse, error) {
	patient := &models.Patient{
		Name:          req.Name,
		Age:           int(req.Age),
		MedicalRecord: req.MedicalRecord,
		BloodType:     req.BloodType,
		Allergies:     req.Allergies,
		Medications:   req.Medications,
	}

	if err := c.service.CreatePatient(patient); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.PatientResponse{
		ID:            uint64(patient.ID),
		Name:          patient.Name,
		Age:           int32(patient.Age),
		MedicalRecord: patient.MedicalRecord,
		BloodType:     patient.BloodType,
		Allergies:     patient.Allergies,
		Medications:   patient.Medications,
	}, nil
}

// Аналогичные реализации для GetPatient, UpdatePatient, DeletePatient
