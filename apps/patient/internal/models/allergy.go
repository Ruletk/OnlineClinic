package models

import (
	"github.com/google/uuid"
	"time"
)

type Allergy struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	PatientID  uuid.UUID `gorm:"type:uuid;not null;index"` // Внешний ключ
	Name       string    `gorm:"not null"`                 // Например, "Пенициллин"
	Severity   string    `gorm:"type:varchar(20)"`         // "LOW", "MODERATE", "SEVERE"
	ObservedAt time.Time `gorm:"type:date"`                // Дата выявления
}
