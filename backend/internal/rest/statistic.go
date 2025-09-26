package rest

import (
	"net/http"
	"prediction_service/models"
	"time"

	"github.com/labstack/echo/v4"
)

type StatisticService interface {
	GetStatistic(start, end time.Time) (*models.StatisticResponse, error)
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
// @Description  Возвращает статистику за указанный период. Если параметры start или end не указаны, используются значения по умолчанию.
// @Tags         statistics
// @Accept       json
// @Produce      json
// @Param        start  query     string  false  "Дата начала периода в формате DD.MM.YYYY"  example("01.01.2024")
// @Param        end    query     string  false  "Дата конца периода в формате DD.MM.YYYY"    example("31.05.2025")
// @Success      200    {object}  models.StatisticResponse
// @Failure      400    {object}  map[string]string "Некорректный формат даты"
// @Failure      500    {object}  map[string]string "Внутренняя ошибка сервиса"
// @Router       /statistics [get]
func (sh *StatisticHandler) StatisticsHandler(c echo.Context) error {
	start := c.QueryParam("start")
	end := c.QueryParam("end")

	layout := "02.01.2006"

	if len(start) == 0 {
		start = "01.01.2024"
	}

	if len(end) == 0 {
		end = "31.05.2025"
	}

	startTime, err := time.Parse(layout, start)
	if err != nil {
		return err
	}

	endTime, err := time.Parse(layout, end)
	if err != nil {
		return err
	}

	resp, err := sh.ss.GetStatistic(startTime, endTime)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
