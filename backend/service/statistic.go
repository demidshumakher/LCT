package service

import (
	"errors"
	"prediction_service/models"
	"time"
)

type StatisticService struct {
}

func NewStatisticService() *StatisticService {
	return &StatisticService{}
}

func (ss *StatisticService) GetStatistic(start, end time.Time) (*models.StatisticResponse, error) {
	return nil, errors.New("Not implemented")
}
