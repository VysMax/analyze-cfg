package usecase

import (
	"github.com/VysMax/analyze-cfg/models"
)

type Usecase interface {
	AnalyzeCfg(cfg *models.Config) ([]models.Problem, error)
}

type Service struct {
	analysis Usecase
}

func New(analysis Usecase) *Service {
	return &Service{analysis: analysis}
}

func (s *Service) AnalyzeCfg(cfg *models.Config) ([]models.Problem, error) {
	problems, err := s.analysis.AnalyzeCfg(cfg)
	if err != nil {
		return nil, err
	}
	return problems, nil
}
