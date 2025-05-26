package service

import (
	"time"

	"github.com/google/uuid"
)

// CreateDoctorRequest — входящие данные для создания
type CreateDoctorRequest struct {
	FirstName        string    `json:"first_name" validate:"required"`
	LastName         string    `json:"last_name"  validate:"required"`
	Patronymic       *string   `json:"patronymic"`
	DateOfBirth      time.Time `json:"date_of_birth" validate:"required"`
	SpecializationID uuid.UUID `json:"specialization_id" validate:"required"`
}

// CreateDoctorResponse — возвращаем клиенту после создания
type CreateDoctorResponse struct {
	ID               uuid.UUID `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Patronymic       *string   `json:"patronymic"`
	DateOfBirth      time.Time `json:"date_of_birth"`
	SpecializationID uuid.UUID `json:"specialization_id"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
}

// DoctorDTO — общий формат отдачи врача
type DoctorDTO struct {
	ID               uuid.UUID `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Patronymic       *string   `json:"patronymic"`
	DateOfBirth      time.Time `json:"date_of_birth"`
	SpecializationID uuid.UUID `json:"specialization_id"`
	Status           string    `json:"status"`
}

// UpdateDoctorRequest — для обновления
type UpdateDoctorRequest struct {
	ID               uuid.UUID `json:"id" validate:"required"`
	FirstName        string    `json:"first_name" validate:"required"`
	LastName         string    `json:"last_name"  validate:"required"`
	Patronymic       *string   `json:"patronymic"`
	DateOfBirth      time.Time `json:"date_of_birth" validate:"required"`
	SpecializationID uuid.UUID `json:"specialization_id" validate:"required"`
	Status           string    `json:"status" validate:"oneof=ACTIVE ON_LEAVE INACTIVE"`
}

// DeleteDoctorResponse — ответ на удаление
type DeleteDoctorResponse struct {
	Success bool `json:"success"`
}
