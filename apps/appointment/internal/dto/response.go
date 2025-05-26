package dto

import (
	"appointment/internal/model"
	"time"
)

type AppointmentResponse struct {
	ID        string                  `json:"id"`
	UserID    string                  `json:"user_id"`
	DoctorID  string                  `json:"doctor_id"`
	Date      string                  `json:"date"`
	Status    model.AppointmentStatus `json:"status"`
	Notes     string                  `json:"notes,omitempty"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}

type AppointmentListResponse struct {
	TotalCount   int                   `json:"total_count"`
	Appointments []AppointmentResponse `json:"appointments"`
}
