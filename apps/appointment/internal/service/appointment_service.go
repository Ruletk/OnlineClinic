package service

import (
	"appointment/internal/dto"
	"appointment/internal/model"
	"appointment/internal/proto/gen"
	"appointment/internal/repository"
	"context"
	"fmt"
	"github.com/Ruletk/OnlineClinic/pkg/logging"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

type AppointmentService interface {
	Create(req *dto.CreateAppointmentRequest) (*dto.AppointmentResponse, error)
	Delete(req *dto.AppointmentIDRequest) error

	GetByUserID(req *dto.GetAppointmentsByUserIDRequest) (*dto.AppointmentListResponse, error)
	GetByID(req *dto.AppointmentIDRequest) (*dto.AppointmentResponse, error)
	GetByDoctorID(req *dto.GetAppointmentsByDoctorIDRequest) (*dto.AppointmentListResponse, error)

	ChangeStatus(req *dto.ChangeAppointmentStatusRequest) (*dto.AppointmentResponse, error)
}

type appointmentService struct {
	repo repository.AppointmentRepository
	grpc repository.GRPCAppointmentRepository
}

func (a appointmentService) Create(req *dto.CreateAppointmentRequest) (*dto.AppointmentResponse, error) {
	logging.Logger.Infof("Creating appointment for user ID: %s with doctor ID: %s on date: %s", req.UserID, req.DoctorID, req.Date)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	availability, err := a.grpc.CheckTimeAvailability(ctx, &gen.CheckTimeAvailabilityRequest{
		DoctorId: req.DoctorID.String(),
		SlotTime: timestamppb.New(req.Date),
	})
	if err != nil {
		logging.Logger.Errorf("Failed to check time availability for doctor ID: %s on date: %s, error: %v", req.DoctorID, req.Date, err)
		return nil, err
	}
	if !availability.IsAvailable {
		logging.Logger.Errorf("Time slot is not available for doctor ID: %s on date: %s with reason: %s", req.DoctorID, req.Date, availability.Reason)
		return nil, fmt.Errorf("time slot is not available for doctor ID: %s on date: %s with reason: %s", req.DoctorID, req.Date, availability.Reason)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = a.changeTimeSlotAvailability(ctx, req.DoctorID, req.Date, false)
	if err != nil {
		logging.Logger.Errorf("Failed to change time slot for doctor ID: %s on date: %s, error: %v", req.DoctorID, req.Date, err)
		return nil, fmt.Errorf("failed to change time slot: %v", err)
	}

	appointment := &model.Appointment{
		UserID:   req.UserID,
		DoctorID: req.DoctorID,
		Date:     req.Date,
		Status:   model.Scheduled,
		Notes:    req.Notes,
	}
	if err := a.repo.Create(appointment); err != nil {
		logging.Logger.Errorf("Failed to create appointment for user ID: %s with doctor ID: %s on date: %s, error: %v", req.UserID, req.DoctorID, req.Date, err)
		return nil, fmt.Errorf("failed to create appointment: %v", err)
	}
	return dto.AppointmentResponseFromModel(appointment), nil
}

func (a appointmentService) Delete(req *dto.AppointmentIDRequest) error {
	logging.Logger.Infof("Deleting appointment with ID: %s", req.ID)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	appointment, err := a.getAppointment(req.ID)
	if err != nil {
		logging.Logger.Errorf("Failed to get appointment by ID: %s, error: %v", req.ID, err)
		return fmt.Errorf("failed to get appointment: %w", err)
	}

	err = a.changeTimeSlotAvailability(ctx, appointment.DoctorID, appointment.Date, true)
	if err != nil {
		logging.Logger.Errorf("Failed to change time slot for doctor ID: %s on date: %s, error: %v", appointment.DoctorID, appointment.Date, err)
		return fmt.Errorf("failed to change time slot: %w", err)
	}

	logging.Logger.Infof("Successfully changed time slot for doctor ID: %s on date: %s", appointment.DoctorID, appointment.Date)
	if err := a.deleteAppointmentRecord(ctx, req.ID, appointment); err != nil {
		logging.Logger.Errorf("Failed to delete appointment with ID: %s, error: %v", req.ID, err)
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	logging.Logger.Infof("Successfully deleted appointment with ID: %s", req.ID)
	return nil
}

func (a appointmentService) GetByUserID(req *dto.GetAppointmentsByUserIDRequest) (*dto.AppointmentListResponse, error) {
	logging.Logger.Infof("Getting appointments for user ID: %s", req.UserID)
	appointments, err := a.repo.ListByUserID(req.UserID)
	if err != nil {
		logging.Logger.Errorf("Failed to get appointments for user ID: %s, error: %v", req.UserID, err)
		return nil, fmt.Errorf("repository get by user ID failed: %w", err)
	}
	if len(appointments) == 0 {
		logging.Logger.Warnf("No appointments found for user ID: %s", req.UserID)
		return &dto.AppointmentListResponse{Appointments: []dto.AppointmentResponse{}}, nil
	}
	return dto.AppointmentListResponseFromModel(appointments), nil
}

func (a appointmentService) GetByID(req *dto.AppointmentIDRequest) (*dto.AppointmentResponse, error) {
	logging.Logger.Infof("Getting appointment by ID: %s", req.ID)
	appointment, err := a.getAppointment(req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment by ID: %w", err)
	}
	return dto.AppointmentResponseFromModel(appointment), nil
}

func (a appointmentService) GetByDoctorID(req *dto.GetAppointmentsByDoctorIDRequest) (*dto.AppointmentListResponse, error) {
	logging.Logger.Infof("Getting appointments for doctor ID: %s", req.DoctorID)
	appointments, err := a.repo.ListByDoctorID(req.DoctorID)
	if err != nil {
		logging.Logger.Errorf("Failed to get appointments for doctor ID: %s, error: %v", req.DoctorID, err)
		return nil, fmt.Errorf("repository get by doctor ID failed: %w", err)
	}
	return dto.AppointmentListResponseFromModel(appointments), nil
}

func (a appointmentService) ChangeStatus(req *dto.ChangeAppointmentStatusRequest) (*dto.AppointmentResponse, error) {
	appointment, err := a.getAppointment(req.ID)
	if err != nil {
		logging.Logger.Errorf("Failed to get appointment by ID: %s, error: %v", req.ID, err)
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	logging.Logger.Infof("Changing status of appointment ID: %s from %s to %s", req.ID, appointment.Status, req.Status)

	appointment.Status = req.Status
	if err := a.repo.Update(appointment); err != nil {
		logging.Logger.Errorf("Failed to update appointment status for ID: %s, error: %v", req.ID, err)
		return nil, fmt.Errorf("failed to update appointment status: %w", err)
	}
	logging.Logger.Infof("Successfully changed status of appointment ID: %s to %s", req.ID, req.Status)
	return dto.AppointmentResponseFromModel(appointment), nil
}

func (a appointmentService) getAppointment(id uuid.UUID) (*model.Appointment, error) {
	appointment, err := a.repo.GetByID(id)
	if err != nil {
		logging.Logger.Errorf("Failed to get appointment by ID: %s, error: %v", id, err)
		return nil, fmt.Errorf("repository get failed: %w", err)
	}
	if appointment == nil {
		logging.Logger.Warnf("Appointment with ID: %s not found", id)
		return nil, gorm.ErrRecordNotFound
	}
	return appointment, nil
}

func (a appointmentService) changeTimeSlotAvailability(ctx context.Context, doctorID uuid.UUID, date time.Time, isAvailable bool) error {
	response, err := a.grpc.ChangeTimeSlot(ctx, &gen.ChangeTimeSlotRequest{
		DoctorId:    doctorID.String(),
		SlotTime:    timestamppb.New(date),
		IsAvailable: isAvailable,
	})
	if err != nil {
		logging.Logger.Errorf("Failed to change time slot for doctor ID: %s on date: %s, error: %v", doctorID, date, err)
		return fmt.Errorf("grpc change time slot failed: %w", err)
	}
	if !response.Success {
		logging.Logger.Errorf("Failed to change time slot for doctor ID: %s on date: %s", doctorID, date)
		return fmt.Errorf("failed to change time slot for doctor ID: %s on date: %s", doctorID, date)
	}
	logging.Logger.Infof("Successfully changed time slot availability to %t for doctor ID: %s on date: %s", isAvailable, doctorID, date)
	return nil
}

func (a appointmentService) deleteAppointmentRecord(ctx context.Context, id uuid.UUID, appointment *model.Appointment) error {
	if err := a.repo.Delete(id); err != nil {
		logging.Logger.Errorf("Failed to delete appointment with ID: %s, error: %v, attempting to revert time slot", id, err)

		if revertErr := a.changeTimeSlotAvailability(ctx, appointment.DoctorID, appointment.Date, false); revertErr != nil {
			logging.Logger.Errorf("Failed to revert time slot for doctor ID: %s on date: %s, error: %v",
				appointment.DoctorID, appointment.Date, revertErr)
			return fmt.Errorf("delete failed: %v, revert also failed: %w", err, revertErr)
		}
		return fmt.Errorf("repository delete failed: %w", err)
	}
	return nil
}

func NewAppointmentService(repo repository.AppointmentRepository) AppointmentService {
	return &appointmentService{
		repo: repo,
	}
}
