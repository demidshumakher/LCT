package postgres

import (
	"database/sql"
	"prediction_service/models"
	"time"
)

type PostgresStatisticRepository struct {
	db *sql.DB
}

func NewPostgresStatisticRepository(db *sql.DB) *PostgresStatisticRepository {
	return &PostgresStatisticRepository{db: db}
}

func (r *PostgresStatisticRepository) GetStatistic(start, end time.Time) (*models.StatisticResponse, error) {
	query := `
	WITH months AS (
		SELECT generate_series(date_trunc('month', $1::date), date_trunc('month', $2::date), interval '1 month') AS month
	)
	SELECT 
		p.name,
		m.month,
		COALESCE(SUM(CASE WHEN r.rating = 'положительно' THEN 1 ELSE 0 END), 0) AS positive,
		COALESCE(SUM(CASE WHEN r.rating = 'негативно' THEN 1 ELSE 0 END), 0) AS negative,
		COALESCE(SUM(CASE WHEN r.rating = 'нейтрально' THEN 1 ELSE 0 END), 0) AS neutral
	FROM Products p
	CROSS JOIN months m
	LEFT JOIN ProductReviews r
		ON r.product_id = p.id 
		AND date_trunc('month', r.date) = m.month
	GROUP BY p.name, m.month
	ORDER BY p.name, m.month;
	`

	rows, err := r.db.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productMap := make(map[string][]models.TimePoint)

	for rows.Next() {
		var name string
		var month time.Time
		var positive, negative, neutral int

		if err := rows.Scan(&name, &month, &positive, &negative, &neutral); err != nil {
			return nil, err
		}

		productMap[name] = append(productMap[name], models.TimePoint{
			Date:     month,
			Positive: positive,
			Negative: negative,
			Neutral:  neutral,
		})
	}

	var response models.StatisticResponse
	for name, timeline := range productMap {
		response.Products = append(response.Products, models.Product{
			Name:     name,
			TimeLine: timeline,
		})
	}

	return &response, nil
}
