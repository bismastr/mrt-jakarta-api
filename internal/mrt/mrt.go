package mrt

import "github.com/bismastr/scrapper-example/internal/repository"

type MrtService struct {
	repo *repository.Queries
}

func NewMrtService(repo *repository.Queries) *MrtService {
	return &MrtService{
		repo: repo,
	}
}

func (s *MrtService) GetAllStation() {

}
