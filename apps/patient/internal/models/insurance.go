package models

import (
	"github.com/google/uuid"
	"time"
)

type Insurance struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	PatientID      uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider       string    `gorm:"not null"`        // Название страховой
	PolicyNumber   string    `gorm:"unique;not null"` // Номер полиса
	ExpirationDate time.Time `gorm:"type:date"`       // Срок действия
}
