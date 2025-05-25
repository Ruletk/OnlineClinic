package services

type AllergyService interface {
}

type allergyService struct {
}

func NewAllergyService() AllergyService {
	return &allergyService{}
}
