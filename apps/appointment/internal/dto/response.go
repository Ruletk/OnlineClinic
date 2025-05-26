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

func AppointmentResponseFromModel(a *model.Appointment) *AppointmentResponse {
	return &AppointmentResponse{
		ID:        a.ID.String(),
		UserID:    a.UserID,
		DoctorID:  a.DoctorID.String(),
		Date:      a.Date.Format(time.RFC3339),
		Status:    a.Status,
		Notes:     a.Notes,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func AppointmentListResponseFromModel(appointments []*model.Appointment) *AppointmentListResponse {
	response := &AppointmentListResponse{
		TotalCount:   len(appointments),
		Appointments: make([]AppointmentResponse, len(appointments)),
	}

	for i, appointment := range appointments {
		response.Appointments[i] = *AppointmentResponseFromModel(appointment)
	}

	return response
}
