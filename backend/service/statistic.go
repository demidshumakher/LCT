package service

import (
	"prediction_service/models"
	"time"
)

type StatisticRepository interface {
	GetStatistic(start, end time.Time) (*models.StatisticResponse, error)
}

type StatisticService struct {
	sr StatisticRepository
}

func NewStatisticService(sr StatisticRepository) *StatisticService {
	return &StatisticService{
		sr: sr,
	}
}

func (ss *StatisticService) GetStatistic(start, end time.Time) (*models.StatisticResponse, error) {
	return ss.sr.GetStatistic(start, end)
}
