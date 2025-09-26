package service

import (
	"prediction_service/models"
)

type StatisticRepository interface {
	GetStatistic() (*models.StatisticResponse, error)
}

type StatisticService struct {
	sr StatisticRepository
}

func NewStatisticService(sr StatisticRepository) *StatisticService {
	return &StatisticService{
		sr: sr,
	}
}

func (ss *StatisticService) GetStatistic() (*models.StatisticResponse, error) {
	return ss.sr.GetStatistic()
}
