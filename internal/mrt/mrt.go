package mrt

import (
	"context"
	"fmt"
	"log"

	"github.com/bismastr/scrapper-example/internal/repository"
)

type MrtService struct {
	repo *repository.Queries
}

func NewMrtService(repo *repository.Queries) *MrtService {
	return &MrtService{
		repo: repo,
	}
}

func (s *MrtService) GetAllStation(ctx context.Context) {
	schedules, err := s.repo.GetSchedule(ctx)
	if err != nil {
		log.Printf("Error getting schedule: %s", err)
	}

	for _, v := range schedules {
		fmt.Println(v.LinesID)
	}
}
