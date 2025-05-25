package services

type InsuranceService interface {
}

type insuranceService struct {
}

func NewInsuranceService() InsuranceService {
	return &insuranceService{}
}
