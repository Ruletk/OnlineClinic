package dto

import (
	"appointment/internal/model"
	"github.com/google/uuid"
	"time"
)

type CreateAppointmentRequest struct {
	UserID   int64     `json:"user_id" binding:"required"`
	DoctorID uuid.UUID `json:"doctor_id" binding:"required"`
	Date     time.Time `json:"date" binding:"required"`
	Notes    string    `json:"notes,omitempty"`
}

type AppointmentIDRequest struct {
	ID uuid.UUID `json:"id" binding:"required,uuid"`
}

type GetAppointmentsByUserIDRequest struct {
	UserID int64 `json:"user_id" binding:"required,uuid"`
}

type GetAppointmentsByDoctorIDRequest struct {
	DoctorID uuid.UUID `json:"doctor_id" binding:"required,uuid"`
}

type ChangeAppointmentStatusRequest struct {
	ID     uuid.UUID               `json:"id" binding:"required,uuid"`
	Status model.AppointmentStatus `json:"status" binding:"required,oneof=scheduled completed canceled"`
}
