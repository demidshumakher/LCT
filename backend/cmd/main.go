package main

import (
	"prediction_service/internal/rest"
	"prediction_service/service"

	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	_ "prediction_service/docs"
)

// @title           API
// @version         1.0
// @description     Server for LCT hackathon
// @host            localhost:8080
// @BasePath        /
func main() {
	e := echo.New()

	ps := service.NewPredictionService()
	rest.NewPredictionHandler(e, ps)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Start(":8080")
}
