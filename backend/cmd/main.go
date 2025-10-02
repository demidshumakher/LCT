package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"prediction_service/internal/repository/postgres"
	"prediction_service/internal/rest"
	"prediction_service/service"

	_ "prediction_service/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// @title           API
// @version         1.0
// @description     Server for LCT hackathon
// @host            hackathon.gberdyshev.org:8080
// @BasePath        /
func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "cryptodb")
	sslMode := getEnv("DB_SSL_MODE", "disable")
	modelUrl := getEnv("MODEL_URL", "")

	// Формирование строки подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORS())

	ss := service.NewStatisticService(postgres.NewPostgresStatisticRepository(db))
	rest.NewStatisticHandler(e, ss)

	predictionCfg := &rest.PredictionConfig{
		Url: modelUrl,
	}
	rest.NewPredictionHandler(e, *predictionCfg)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Start("0.0.0.0:8080")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
