package repository

import (
	"appointment/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	appointmentCacheTTL = 5 * time.Minute
)

type AppointmentDBRedisRepository struct {
	repo  AppointmentRepository
	redis RedisRepository
	ctx   context.Context
}

func (a *AppointmentDBRedisRepository) getCacheKey(id uuid.UUID) string {
	return fmt.Sprintf("appointment:%s", id.String())
}

func (a *AppointmentDBRedisRepository) getListCacheKey(userID int64, doctorID uuid.UUID) string {
	if userID != 0 {
		return fmt.Sprintf("appointments:user:%d", userID)
	}
	return fmt.Sprintf("appointments:doctor:%s", doctorID.String())
}

func (a *AppointmentDBRedisRepository) Create(appointment *model.Appointment) error {
	// Create in DB first
	if err := a.repo.Create(appointment); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	_ = a.redis.Delete(ctx, a.getListCacheKey(appointment.UserID, uuid.Nil))
	_ = a.redis.Delete(ctx, a.getListCacheKey(0, appointment.DoctorID))

	return nil
}

func (a *AppointmentDBRedisRepository) GetByID(id uuid.UUID) (*model.Appointment, error) {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	cacheKey := a.getCacheKey(id)

	// Try to get from cache first
	cached, err := a.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var appointment model.Appointment
		if err := json.Unmarshal([]byte(cached), &appointment); err == nil {
			return &appointment, nil
		}
	}

	// Cache miss or error - get from DB
	appointment, err := a.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the result
	appointmentBytes, err := json.Marshal(appointment)
	if err == nil {
		_ = a.redis.Set(ctx, cacheKey, string(appointmentBytes), appointmentCacheTTL)
	}

	return appointment, nil
}

func (a *AppointmentDBRedisRepository) Update(appointment *model.Appointment) error {
	// Update in DB first
	if err := a.repo.Update(appointment); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	cacheKey := a.getCacheKey(appointment.ID)

	// Invalidate the cached appointment
	_ = a.redis.Delete(ctx, cacheKey)

	// Also invalidate any lists that might contain this appointment
	_ = a.redis.Delete(ctx, a.getListCacheKey(appointment.UserID, uuid.Nil))
	_ = a.redis.Delete(ctx, a.getListCacheKey(0, appointment.DoctorID))

	return nil
}

func (a *AppointmentDBRedisRepository) Delete(id uuid.UUID) error {
	// First get the appointment to know which caches to invalidate
	appointment, err := a.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete from DB
	if err := a.repo.Delete(id); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	cacheKey := a.getCacheKey(id)

	// Invalidate the cached appointment
	_ = a.redis.Delete(ctx, cacheKey)

	// Invalidate any lists that might contain this appointment
	if appointment != nil {
		_ = a.redis.Delete(ctx, a.getListCacheKey(appointment.UserID, uuid.Nil))
		_ = a.redis.Delete(ctx, a.getListCacheKey(0, appointment.DoctorID))
	}

	return nil
}

func (a *AppointmentDBRedisRepository) ListByUserID(userID int64) ([]model.Appointment, error) {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	cacheKey := a.getListCacheKey(userID, uuid.Nil)

	// Try to get from cache first
	cached, err := a.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var appointments []model.Appointment
		if err := json.Unmarshal([]byte(cached), &appointments); err == nil {
			return appointments, nil
		}
	}

	// Cache miss or error - get from DB
	appointments, err := a.repo.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	appointmentsBytes, err := json.Marshal(appointments)
	if err == nil {
		_ = a.redis.Set(ctx, cacheKey, string(appointmentsBytes), appointmentCacheTTL)
	}

	return appointments, nil
}

func (a *AppointmentDBRedisRepository) ListByDoctorID(doctorID uuid.UUID) ([]model.Appointment, error) {
	ctx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	cacheKey := a.getListCacheKey(0, doctorID)

	// Try to get from cache first
	cached, err := a.redis.Get(ctx, cacheKey)
	if err == nil && cached != "" {
		var appointments []model.Appointment
		if err := json.Unmarshal([]byte(cached), &appointments); err == nil {
			return appointments, nil
		}
	}

	// Cache miss or error - get from DB
	appointments, err := a.repo.ListByDoctorID(doctorID)
	if err != nil {
		return nil, err
	}

	// Cache the result
	appointmentsBytes, err := json.Marshal(appointments)
	if err == nil {
		_ = a.redis.Set(ctx, cacheKey, string(appointmentsBytes), appointmentCacheTTL)
	}

	return appointments, nil
}

func NewAppointmentDBRedisRepository(repo AppointmentRepository, redis RedisRepository, ctx context.Context) AppointmentRepository {
	return &AppointmentDBRedisRepository{
		repo:  repo,
		redis: redis,
		ctx:   ctx,
	}
}
