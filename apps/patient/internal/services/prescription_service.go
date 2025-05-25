package services

type PrescriptionService interface {
}

type prescriptionService struct {
}

func NewPrescriptionService() PrescriptionService {
	return &prescriptionService{}

}
