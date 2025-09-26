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
	"github.com/swaggo/echo-swagger"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// @title           API
// @version         1.0
// @description     Server for LCT hackathon
// @host            localhost:8080
// @BasePath        /
func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "cryptodb")
	sslMode := getEnv("DB_SSL_MODE", "disable")

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

	ss := service.NewStatisticService(postgres.NewPostgresStatisticRepository(db))
	rest.NewStatisticHandler(e, ss)

	ps := service.NewPredictionService()
	rest.NewPredictionHandler(e, ps)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	testData(db)

	e.Start("0.0.0.0:8080")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func testData(db *sql.DB) {
	db.Exec(q)
}

const q = `
		INSERT INTO Products (name) VALUES
('Смартфон Xiaomi Redmi Note 12'),
('Ноутбук ASUS VivoBook 15'),
('Наушники Sony WH-1000XM4'),
('Фитнес-браслет Huawei Band 8'),
('Игровая мышь Logitech G Pro X');

-- Вставляем отзывы (100 записей)
INSERT INTO ProductReviews (product_id, date, rating) VALUES
(1, '2024-01-15', 'положительно'),
(1, '2024-01-18', 'положительно'),
(1, '2024-02-02', 'нейтрально'),
(1, '2024-02-14', 'негативно'),
(1, '2024-03-01', 'положительно'),
(2, '2024-01-20', 'положительно'),
(2, '2024-02-05', 'положительно'),
(2, '2024-02-22', 'нейтрально'),
(2, '2024-03-10', 'негативно'),
(2, '2024-03-12', 'положительно'),
(3, '2024-01-25', 'положительно'),
(3, '2024-02-10', 'положительно'),
(3, '2024-02-28', 'нейтрально'),
(3, '2024-03-05', 'негативно'),
(3, '2024-03-18', 'положительно'),
(4, '2024-02-01', 'положительно'),
(4, '2024-02-16', 'положительно'),
(4, '2024-03-02', 'нейтрально'),
(4, '2024-03-08', 'негативно'),
(4, '2024-03-20', 'положительно'),
(5, '2024-02-12', 'положительно'),
(5, '2024-02-25', 'положительно'),
(5, '2024-03-07', 'нейтрально'),
(5, '2024-03-14', 'негативно'),
(5, '2024-03-14', 'негативно'),
(5, '2024-03-14', 'негативно'),
(5, '2024-03-14', 'негативно'),
(5, '2024-03-14', 'негативно'),
(5, '2024-03-22', 'положительно'),
(5, '2024-03-22', 'положительно'),
(5, '2024-03-22', 'положительно'),
(5, '2024-03-22', 'негативно'),
(5, '2024-03-22', 'положительно');
`
