package repository

import (
	"appointment/internal/proto/gen"
	"context"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"time"
)

type GRPCAppointmentRepository interface {
	CheckTimeAvailability(ctx context.Context, req *gen.CheckTimeAvailabilityRequest) (*gen.CheckTimeAvailabilityResponse, error)
	GetAvailableSlots(ctx context.Context, req *gen.GetAvailableSlotsRequest) (*gen.GetAvailableSlotsResponse, error)
	ChangeTimeSlot(ctx context.Context, req *gen.ChangeTimeSlotRequest) (*gen.ChangeTimeSlotResponse, error)
}

type grpcAppointmentRepository struct {
	client         gen.DoctorServiceClient
	defaultTimeout time.Duration
}

func NewGRPCAppointmentRepository(client gen.DoctorServiceClient) GRPCAppointmentRepository {
	logging.Logger.Info("Initializing gRPC Appointment Repository")
	return &grpcAppointmentRepository{
		client: client,
		// TODO: Warning, default timeout is hardcoded, consider making it configurable
		defaultTimeout: 5 * time.Second, // Default timeout for gRPC calls
	}
}

func (g *grpcAppointmentRepository) CheckTimeAvailability(ctx context.Context, req *gen.CheckTimeAvailabilityRequest) (*gen.CheckTimeAvailabilityResponse, error) {
	logging.Logger.Infof("Checking time availability for doctor ID: %s, date: %s", req.DoctorId, req.SlotTime)
	if err := g.checkClient(); err != nil {
		return nil, err
	}
	c, cancel := context.WithTimeout(ctx, g.defaultTimeout)
	defer cancel()

	return g.client.CheckTimeAvailability(c, req)
}

func (g *grpcAppointmentRepository) GetAvailableSlots(ctx context.Context, req *gen.GetAvailableSlotsRequest) (*gen.GetAvailableSlotsResponse, error) {
	logging.Logger.Infof("Getting available slots for doctor ID: %s, and date between: %s and %s", req.DoctorId, req.StartDate, req.EndDate)
	if err := g.checkClient(); err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, g.defaultTimeout)
	defer cancel()

	return g.client.GetAvailableSlots(c, req)
}

func (g *grpcAppointmentRepository) ChangeTimeSlot(ctx context.Context, req *gen.ChangeTimeSlotRequest) (*gen.ChangeTimeSlotResponse, error) {
	logging.Logger.Infof("Changing time slot for doctor: %s, time: %s, status: %v", req.DoctorId, req.SlotTime, req.IsAvailable)
	if err := g.checkClient(); err != nil {
		return nil, err
	}

	c, cancel := context.WithTimeout(ctx, g.defaultTimeout)
	defer cancel()

	return g.client.ChangeTimeSlot(c, req)
}

func (g *grpcAppointmentRepository) checkClient() error {
	if g.client == nil {
		logging.Logger.Error("g.client is nil, ensure the grpcAppointmentRepository is properly initialized")
		return fmt.Errorf("g.client is nil, ensure the grpcAppointmentRepository is properly initialized")
	}
	return nil
}
