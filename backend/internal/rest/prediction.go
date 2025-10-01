package rest

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PredictionConfig struct {
	Url string
}

type PredictionHandler struct {
	cfg PredictionConfig
}

func NewPredictionHandler(e *echo.Echo, cfg PredictionConfig) {
	ph := &PredictionHandler{
		cfg: cfg,
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
	req, err := http.NewRequest(http.MethodPost, ph.cfg.Url, c.Request().Body)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return c.Blob(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}
