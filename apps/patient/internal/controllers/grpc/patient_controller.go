package grpcpackage

import (
	proto "patient/internal/proto/gen"
	"patient/internal/services"
)

type PatientController struct {
	proto.UnimplementedPatientServiceServer
	service services.PatientService
}

func NewPatientController(svc services.PatientService) *PatientController {
	return &PatientController{service: svc}
}

//func (c *PatientController) CreatePatient(req *proto.PatientRequest) (*proto.PatientResponse, error) {
//	patient := &models.Patient{
//		Height:    req.Height,
//		BloodType: req.BloodType,
//		Weight:    req.Weight,
//	}
//
//	if err := c.service.CreatePatient(patient); err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &proto.PatientResponse{
//		ID:         uint64(patient.ID),
//		BloodType:  patient.BloodType,
//		Height:     patient.Height,
//		Weight:     patient.Weight,
//		Allergies:  patient.Allergies,
//		Insurances: patient.Insurances,
//	}, nil
//}

// Аналогичные реализации для GetPatient, UpdatePatient, DeletePatient
