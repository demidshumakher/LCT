package models

import (
	"encoding/json"
	"time"
)

type StatisticResponse struct {
	Products []Product `json:"products"`
}

type Product struct {
	Name     string      `json:"name" example:"ипотека"`
	TimeLine []TimePoint `json:"timeline"`
}

type TimePoint struct {
	Date     time.Time `json:"date" example:"01.01.2025"`
	Positive int       `json:"positive" example:"50"`
	Negative int       `json:"negative" example:"100"`
	Neutral  int       `json:"neutral" example:"10"`
}

func (t TimePoint) MarshalJSON() ([]byte, error) {
	type Alias TimePoint
	return json.Marshal(&struct {
		Date string `json:"date"`
		Alias
	}{
		Date:  t.Date.Format("2006-01"),
		Alias: (Alias)(t),
	})
}
