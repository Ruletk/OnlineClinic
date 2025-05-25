package models

import (
	"github.com/google/uuid"
	"time"
)

type Prescription struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	PatientID  uuid.UUID `gorm:"type:uuid;not null;index"` // Внешний ключ
	DoctorID   uuid.UUID `gorm:"type:uuid;not null"`       // Ссылка на Doctor сервис
	Medication string    `gorm:"not null"`                 // Название препарата
	Dosage     string    `gorm:"not null"`                 // "500 мг 2 раза в день"
	ValidUntil time.Time `gorm:"type:date"`                // Дата окончания
}
