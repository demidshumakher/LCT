package rest

import (
	"net/http"
	"prediction_service/models"

	"github.com/labstack/echo/v4"
)

type StatisticService interface {
	GetStatistic() (*models.StatisticResponse, error)
}

type StatisticHandler struct {
	ss StatisticService
}

func NewStatisticHandler(e *echo.Echo, ss StatisticService) {
	sh := &StatisticHandler{
		ss: ss,
	}

	e.GET("/statistics", sh.StatisticsHandler)
}

// StatisticsHandler godoc
// @Summary      Получение статистики
// @Description  Возвращает статистику за весь период.
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Success      200    {object}  models.StatisticResponse
// @Failure      500    {object}  map[string]string "Внутренняя ошибка сервиса"
// @Router       /statistics [get]
func (sh *StatisticHandler) StatisticsHandler(c echo.Context) error {
	resp, err := sh.ss.GetStatistic()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
