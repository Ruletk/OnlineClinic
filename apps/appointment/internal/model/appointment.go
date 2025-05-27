package model

import (
	"github.com/google/uuid"
	"time"
)

type Appointment struct {
	ID       uuid.UUID         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID   int64             `json:"user_id" gorm:"type:uuid;not null"`
	DoctorID uuid.UUID         `json:"doctor_id" gorm:"type:uuid;not null"`
	Date     time.Time         `json:"date" gorm:"not null"`
	Status   AppointmentStatus `json:"status" gorm:"type:ENUM('scheduled', 'completed', 'canceled');not null;default:'scheduled'"`
	Notes    string            `json:"notes" gorm:"type:text;default:null"` // Optional notes for the appointment

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"` // Soft delete
}

type AppointmentStatus string

const (
	Scheduled AppointmentStatus = "scheduled"
	Completed AppointmentStatus = "completed"
	Canceled  AppointmentStatus = "canceled"
)
