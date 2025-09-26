package service

import (
	"errors"
	"prediction_service/models"
)

type PredictionService struct {
}

func NewPredictionService() *PredictionService {
	return &PredictionService{}
}

func (ps *PredictionService) Predict(req *models.PredictionRequest) (*models.PredictionResponse, error) {
	return nil, errors.New("Not Implemented")
}
