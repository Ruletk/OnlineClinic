package repository

import (
	"appointment/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	Create(appointment *model.Appointment) (*model.Appointment, error)
	GetByID(id uuid.UUID) (*model.Appointment, error)
	Update(appointment *model.Appointment) (*model.Appointment, error)
	Delete(id uuid.UUID) error
	ListByUserID(userID uuid.UUID) ([]model.Appointment, error)
	ListByDoctorID(doctorID uuid.UUID) ([]model.Appointment, error)
}

type appointmentRepository struct {
	db *gorm.DB
}

func (a appointmentRepository) Create(appointment *model.Appointment) (*model.Appointment, error) {
	if err := a.db.Create(appointment).Error; err != nil {
		return nil, err
	}
	return appointment, nil
}

func (a appointmentRepository) GetByID(id uuid.UUID) (*model.Appointment, error) {
	var appointment model.Appointment
	if err := a.db.First(&appointment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (a appointmentRepository) Update(appointment *model.Appointment) (*model.Appointment, error) {
	if err := a.db.Save(appointment).Error; err != nil {
		return nil, err
	}
	return appointment, nil
}

func (a appointmentRepository) Delete(id uuid.UUID) error {
	return a.db.Delete(&model.Appointment{}, "id = ?", id).Error
}

func (a appointmentRepository) ListByUserID(userID uuid.UUID) ([]model.Appointment, error) {
	var appointments []model.Appointment
	if err := a.db.Where("user_id = ?", userID).Find(&appointments).Error; err != nil {
		return nil, err
	}
	return appointments, nil
}

func (a appointmentRepository) ListByDoctorID(doctorID uuid.UUID) ([]model.Appointment, error) {
	var appointments []model.Appointment
	if err := a.db.Where("doctor_id = ?", doctorID).Find(&appointments).Error; err != nil {
		return nil, err
	}
	return appointments, nil
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}
