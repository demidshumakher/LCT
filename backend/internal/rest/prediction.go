package rest

import (
	"net/http"
	"prediction_service/models"

	"github.com/labstack/echo/v4"
)

type PredictionService interface {
	Predict(*models.PredictionRequest) (*models.PredictionResponse, error)
}

type PredictionHandler struct {
	ps PredictionService
}

func NewPredictionHandler(e *echo.Echo, ps PredictionService) {
	ph := &PredictionHandler{
		ps: ps,
	}

	e.POST("/predict", ph.PredictionHandler)
}

// PredictionHandler godoc
// @Summary      Анализ комментариев
// @Description  Принимает список комментариев и возвращает предсказанные темы и тональности для каждого
// @Tags         prediction
// @Accept       json
// @Produce      json
// @Param        request  body      models.PredictionRequest  true  "Комментарии для анализа"
// @Success      200      {object}  models.PredictionResponse
// @Failure      400      {object}  map[string]string
// @Router       /predict [post]
func (ph *PredictionHandler) PredictionHandler(c echo.Context) error {
	req := new(models.PredictionRequest)

	if err := c.Bind(req); err != nil {
		return err
	}

	resp, err := ph.ps.Predict(req)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
