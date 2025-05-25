package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Patient struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID   `gorm:"type:uuid;not null;index"` // Ссылка на User сервис
	BloodType  string      `gorm:"type:varchar(5)"`          // Например: "A+", "O-"
	Height     float64     `gorm:"type:decimal(5,2)"`        // Рост в см
	Weight     float64     `gorm:"type:decimal(5,2)"`        // Вес в кг
	Allergies  []Allergy   `gorm:"foreignKey:PatientID"`     // Часто нужны → eager load
	Insurances []Insurance `gorm:"foreignKey:PatientID"`     // Часто нужны → eager load
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
}
