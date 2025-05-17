package handler

import (
	"context"
	"doctor/internal/model"
	"doctor/internal/usecase"
	pb "genproto"
)

type GrpcHandler struct {
	usecase usecase.UseCase
	pb.UnimplementedDoctorServiceServer
}

func NewGrpcHandler(uc usecase.UseCase) *GrpcHandler {
	return &GrpcHandler{usecase: uc}
}

// CreateDoctor handles creation of a new doctor
func (h *GrpcHandler) CreateDoctor(ctx context.Context, req *pb.CreateDoctorRequest) (*pb.DoctorResponse, error) {
	doc := &model.Doctor{
		Name:      req.Name,
		Specialty: req.Specialty,
		Email:     req.Email,
	}
	err := h.usecase.Create(ctx, doc)
	if err != nil {
		return nil, err
	}

	return &pb.DoctorResponse{
		Id:        doc.ID,
		Name:      doc.Name,
		Specialty: doc.Specialty,
		Email:     doc.Email,
	}, nil
}

// GetDoctorByID returns a doctor by ID
func (h *GrpcHandler) GetDoctorByID(ctx context.Context, req *pb.GetDoctorRequest) (*pb.DoctorResponse, error) {
	doc, err := h.usecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DoctorResponse{
		Id:        doc.ID,
		Name:      doc.Name,
		Specialty: doc.Specialty,
		Email:     doc.Email,
	}, nil
}

// UpdateDoctor updates doctor information
func (h *GrpcHandler) UpdateDoctor(ctx context.Context, req *pb.UpdateDoctorRequest) (*pb.DoctorResponse, error) {
	doc := &model.Doctor{
		ID:        req.Id,
		Name:      req.Name,
		Specialty: req.Specialty,
		Email:     req.Email,
	}
	err := h.usecase.Update(ctx, doc)
	if err != nil {
		return nil, err
	}

	return &pb.DoctorResponse{
		Id:        doc.ID,
		Name:      doc.Name,
		Specialty: doc.Specialty,
		Email:     doc.Email,
	}, nil
}

// DeleteDoctor deletes a doctor by ID
func (h *GrpcHandler) DeleteDoctor(ctx context.Context, req *pb.DeleteDoctorRequest) (*pb.DeleteDoctorResponse, error) {
	err := h.usecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteDoctorResponse{
		Status: "deleted",
	}, nil
}
