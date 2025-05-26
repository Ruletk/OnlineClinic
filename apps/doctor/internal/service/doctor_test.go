package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Ruletk/OnlineClinic/apps/doctor/internal/model"
	"github.com/Ruletk/OnlineClinic/apps/doctor/internal/repository"
)

func TestDoctorService_UpdateDoctor_Success(t *testing.T) {
	repo := new(MockDoctorRepo)
	svc := NewDoctorService(repo, nil)

	// исходная запись
	id := uuid.New()
	orig := &model.Doctor{
		ID:               id,
		FirstName:        "Old",
		LastName:         "Name",
		Patronymic:       "Patr",
		DateOfBirth:      time.Date(1990, 2, 2, 0, 0, 0, 0, time.UTC),
		SpecializationID: uuid.New(),
		Status:           model.Active,
	}

	// запрос на обновление
	newSpec := uuid.New()
	req := UpdateDoctorRequest{
		ID:               id,
		FirstName:        "NewFirst",
		LastName:         "NewLast",
		Patronymic:       "NewPatr",
		DateOfBirth:      time.Date(1991, 3, 3, 0, 0, 0, 0, time.UTC),
		SpecializationID: newSpec,
		Status:           string(model.OnLeave),
	}

	// 1) когда GetByID, вернуть orig
	repo.
		On("GetByID", mock.Anything, id).
		Return(orig, nil)

	// 2) когда Update, проверить, что поля подменились
	repo.
		On("Update", mock.Anything, mock.MatchedBy(func(d *model.Doctor) bool {
			return d.ID == id &&
				d.FirstName == req.FirstName &&
				d.LastName == req.LastName &&
				d.Patronymic == req.Patronymic &&
				d.DateOfBirth.Equal(req.DateOfBirth) &&
				d.SpecializationID == req.SpecializationID &&
				d.Status == model.OnLeave
		})).
		Return(nil)

	dto, err := svc.UpdateDoctor(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, req.ID, dto.ID)
	require.Equal(t, req.FirstName, dto.FirstName)
	require.Equal(t, req.Status, dto.Status)

	repo.AssertExpectations(t)
}

func TestDoctorService_UpdateDoctor_NotFound(t *testing.T) {
	repo := new(MockDoctorRepo)
	svc := NewDoctorService(repo, nil)

	id := uuid.New()
	req := UpdateDoctorRequest{ID: id}

	repo.
		On("GetByID", mock.Anything, id).
		Return((*model.Doctor)(nil), repository.ErrRecordNotFound)

	_, err := svc.UpdateDoctor(context.Background(), req)
	// метод в текущей реализации просто возвращает ErrRecordNotFound
	require.ErrorIs(t, err, repository.ErrRecordNotFound)

	repo.AssertExpectations(t)
}

func TestDoctorService_DeleteDoctor_Success(t *testing.T) {
	repo := new(MockDoctorRepo)
	svc := NewDoctorService(repo, nil)

	id := uuid.New()
	repo.
		On("Delete", mock.Anything, id).
		Return(nil)

	resp, err := svc.DeleteDoctor(context.Background(), id)
	require.NoError(t, err)
	require.True(t, resp.Success)

	repo.AssertExpectations(t)
}

func TestDoctorService_DeleteDoctor_Error(t *testing.T) {
	repo := new(MockDoctorRepo)
	svc := NewDoctorService(repo, nil)

	id := uuid.New()
	repo.
		On("Delete", mock.Anything, id).
		Return(errors.New("delete failed"))

	_, err := svc.DeleteDoctor(context.Background(), id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "delete failed")

	repo.AssertExpectations(t)
}
